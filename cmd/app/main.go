package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/g6834/team17/auth-service/internal/api/grpc"
	"gitlab.com/g6834/team17/auth-service/internal/api/handlers"
	"gitlab.com/g6834/team17/auth-service/internal/api/middlewares"
	"gitlab.com/g6834/team17/auth-service/internal/api/presenters"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"gitlab.com/g6834/team17/auth-service/internal/databases"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/repositories"
	"gitlab.com/g6834/team17/auth-service/internal/services/auth_service"
	"gitlab.com/g6834/team17/auth-service/internal/services/user_service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	configName = "config.yaml"
)

var MigrationNoChange = errors.New("no change")

func main() {
	if err := config.ReadConfigYML(configName); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.New()

	debug := flag.Bool("debug", cfg.App.Debug, "Defines log level to debug")
	useDatabase := flag.Bool("use-db", true, "Defines repository source")
	migration := flag.Bool("migration", true, "Defines the migration start option")
	telemetry := flag.Bool("telemetry", true, "Defines the telemetry start option")
	flag.Parse()

	logger := initLogger(cfg, *debug)

	//*telemetry = false // todo: delete when will be created jaeger in k8s
	if *telemetry {
		initOtel(&cfg, logger)
	}

	// Repositories
	var userRepo interfaces.UserRepo

	*useDatabase = false // todo: delete when will be created mongodb in k8s
	if *useDatabase {
		// Database
		mongo := initMongo(&cfg)

		// Migration
		if *migration {
			runMigrations(mongo, &cfg)
		}

		userRepo = repositories.NewDatabaseRepo(mongo)
	} else {
		userRepo = repositories.NewFileRepo()
	}

	// Presenters
	presenters := presenters.NewPresenters(logger)

	// Services
	authService := auth_service.New(&auth_service.JwtSettings{
		SecretKey:  cfg.Jwt.SecretKey,
		AtLifeTime: cfg.Jwt.AtLifeTime,
		RtLifeTime: cfg.Jwt.RtLifeTime,
	}, userRepo)
	userService := user_service.New(userRepo)

	restRouter := chi.NewRouter()
	restRouter.Route("/v1", func(r chi.Router) {
		r.Use(middleware.RealIP)
		r.Use(middlewares.RequestID)
		r.Use(middlewares.Tracer)
		r.Use(middlewares.Logger(logger))
		r.Use(middlewares.Recover(logger))
		r.Use(cors.Default().Handler)

		r.Mount("/auth", handlers.AuthRouter(logger, presenters, authService))

		r.With(middlewares.Validate(presenters, authService)).
			Mount("/user", handlers.UserRouter(logger, presenters, userService))
	})

	restAddress := fmt.Sprintf("%v:%v", cfg.Rest.Host, cfg.Rest.Port)
	restSrv := http.Server{
		Addr:         restAddress,
		Handler:      restRouter,
		ReadTimeout:  time.Second * time.Duration(cfg.Rest.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Rest.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Rest.IdleTimeout),
	}

	debugRouter := chi.NewRouter()
	debugRouter.Mount("/debug", handlers.ProfilerRouter(logger, presenters))
	debugAddress := fmt.Sprintf("%v:%v", cfg.Rest.Host, cfg.Rest.DebugPort)
	debugSrv := http.Server{
		Addr:    debugAddress,
		Handler: debugRouter,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	// Graceful shutdown
	go func() {
		<-signalChan
		logger.Warn().Msg("shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Rest.ShutdownTimeout)*time.Second)
		defer cancel()
		err := restSrv.Shutdown(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("couldn't terminated rest server")
		}

		err = debugSrv.Shutdown(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("couldn't terminated debug server")
		}

		select {
		case <-time.After(time.Duration(cfg.Rest.ShutdownTimeout+1) * time.Second):
			logger.Warn().Msg("not all connections done")
		case <-ctx.Done():
		}
	}()

	go func() {
		if err := grpc.NewGrpcServer(authService).Start(&cfg); err != nil {
			logger.Fatal().Err(err).Msg("Failed creating gRPC server")
		}
	}()

	go func() {
		if err := debugSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("debug server was terminated with an error")
		}
	}()

	log.Info().
		Str("version", cfg.App.Version).
		Bool("debug", *debug).
		Str("environment", cfg.App.Environment).
		Msgf("Starting services: %s", cfg.App.Name)

	if err := restSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("rest server was terminated with an error")
	}
}

func initLogger(cfg config.Config, debug bool) zerolog.Logger {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("app_name", cfg.App.Name).
		Str("host_ip", cfg.Rest.Host).
		Str("host_port", fmt.Sprint(cfg.Rest.Port)).
		Logger()

	return logger
}

func initOtel(cfg *config.Config, logger zerolog.Logger) {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Jaeger.Host), jaeger.WithAgentPort(cfg.Jaeger.Port)),
	)

	if err != nil {
		logger.Fatal().Err(err).Msg("failed connecting to apm exporter")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Jaeger.Service),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func runMigrations(mongo *mongo.Database, cfg *config.Config) {
	driver, err := mongodb.WithInstance(mongo.Client(), &mongodb.Config{DatabaseName: cfg.Database.Name})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't get databases driver")
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+cfg.Database.Migrations, cfg.Database.Name, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}
	if err := m.Up(); err != nil && errors.Is(err, MigrationNoChange) {
		log.Fatal().Err(err).Msg("Migration failed")
	}
}

func initMongo(cfg *config.Config) *mongo.Database {
	mongoCfg := db.MongoConfig{
		Timeout:  cfg.Database.Timeout * int(time.Second),
		DBname:   cfg.Database.Name,
		Username: cfg.Database.User,
		Password: cfg.Database.Password,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
	}

	mongo, err := db.Connect(mongoCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't connection to MongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = mongo.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't ping to MongoDB")
	}

	return mongo
}
