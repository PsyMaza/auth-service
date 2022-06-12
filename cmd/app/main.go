package main

import (
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
	"gitlab.com/g6834/team17/auth-service/internal/api/handlers"
	"gitlab.com/g6834/team17/auth-service/internal/api/middlewares"
	"gitlab.com/g6834/team17/auth-service/internal/api/presenters"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"gitlab.com/g6834/team17/auth-service/internal/databases"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/repo"
	"gitlab.com/g6834/team17/auth-service/internal/services/auth_service"
	"gitlab.com/g6834/team17/auth-service/internal/services/user_service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"net/http"
	"os"
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

	debug := flag.Bool("debug", cfg.App.Debug, "sets log level to debug")
	useDatabase := flag.Bool("use-db", true, "sets repository source")
	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

	log.Info().
		Str("version", cfg.App.Version).
		Bool("debug", *debug).
		Str("environment", cfg.App.Environment).
		Msgf("Starting services: %s", cfg.App.Name)

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("role", cfg.App.Name).
		Str("host", cfg.Rest.Host).
		Logger()

	initOtel(&cfg, logger)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Database
	mongo := initMongo(&cfg)

	// Repositories
	var userRepo interfaces.UserRepo

	//*useDatabase = false // todo: delete when will be created mongodb in k8s
	if *useDatabase {
		userRepo = repo.NewDatabaseRepo(mongo)
	} else {
		userRepo = nil
	}

	// Presenters
	presenters := presenters.NewPresenters(logger)

	// Migration
	*migration = false // todo: delete when will be created mongodb in k8s
	if *migration {
		runMigrations(mongo, &cfg)
	}

	// Services
	authService := auth_service.New(&auth_service.JwtSettings{
		SecretKey:  cfg.Jwt.SecretKey,
		AtLifeTime: cfg.Jwt.AtLifeTime,
		RtLifeTime: cfg.Jwt.RtLifeTime,
	}, userRepo)
	userService := user_service.New(userRepo)

	router := chi.NewRouter()
	router.Route("/", func(v1 chi.Router) {
		v1.Use(middleware.RealIP)
		v1.Use(middlewares.RequestID)
		v1.Use(middlewares.Tracer)
		v1.Use(middlewares.Logger(logger))
		v1.Use(middlewares.Recover(logger))
		v1.Use(cors.Default().Handler)

		v1.Mount("/auth/v1", handlers.AuthRouter(logger, presenters, authService))
		v1.Mount("/user/v1", handlers.UserRouter(logger, presenters, userService))
	})

	listenAddress := fmt.Sprintf("%v:%v", cfg.Rest.Host, cfg.Rest.Port)
	http.ListenAndServe(listenAddress, router)
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

	return mongo
}
