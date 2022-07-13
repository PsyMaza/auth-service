package application

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/g6834/team17/auth-service/internal/api/grpc"
	"gitlab.com/g6834/team17/auth-service/internal/api/handlers"
	"gitlab.com/g6834/team17/auth-service/internal/api/middlewares"
	"gitlab.com/g6834/team17/auth-service/internal/api/presenters"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"gitlab.com/g6834/team17/auth-service/internal/infrastructure"
	"gitlab.com/g6834/team17/auth-service/internal/repositories"
	"gitlab.com/g6834/team17/auth-service/internal/services/auth_service"
	"gitlab.com/g6834/team17/auth-service/internal/services/user_service"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

const (
	configName = "config.yaml"
)

var (
	servers []*infrastructure.Server
	logger  *zerolog.Logger
)

func Start(ctx context.Context) {
	if err := config.ReadConfigYML(configName); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.New()

	debug := flag.Bool("debug", cfg.App.Debug, "Defines log level to debug")
	flag.Parse()

	logger = infrastructure.NewLogger(cfg, *debug)

	infrastructure.NewTelemetry(cfg, logger)

	// Database
	mongo := infrastructure.NewDatabase(ctx, cfg, logger)

	// Repositories
	userRepo := repositories.NewDatabaseRepo(mongo)

	// Presenters
	presenters := presenters.NewPresenters(logger)

	// Services
	authService := auth_service.New(&auth_service.JwtSettings{
		SecretKey:  cfg.Jwt.SecretKey,
		AtLifeTime: cfg.Jwt.AtLifeTime,
		RtLifeTime: cfg.Jwt.RtLifeTime,
	}, userRepo)
	userService := user_service.New(userRepo)

	var g errgroup.Group

	g.Go(func() error {
		err := grpc.NewGrpcServer(authService).Start(cfg)

		return fmt.Errorf("failed creating grpc server. %w", err)
	})

	g.Go(func() error {
		debugRouter := chi.NewMux()
		debugRouter.Mount("/debug", handlers.ProfilerRouter(logger, presenters))
		debugAddress := fmt.Sprintf("%v:%v", cfg.Http.Host, cfg.Http.DebugPort)
		debugSrv, err := infrastructure.NewServer(logger, debugRouter, debugAddress, cfg)
		if err != nil {
			return fmt.Errorf("debug server creating failed. %w", err)
		}
		servers = append(servers, debugSrv)

		err = debugSrv.Start()

		return fmt.Errorf("debug server was terminated with an error. %w", err)
	})

	g.Go(func() error {
		swaggerAddress := fmt.Sprintf("%v:%v", cfg.Http.Host, cfg.Http.SwaggerPort)
		swaggerRouter := handlers.SwaggerRouter(swaggerAddress)
		swaggerSrv, err := infrastructure.NewServer(logger, swaggerRouter, swaggerAddress, cfg)
		if err != nil {
			return fmt.Errorf("swagger server creating failed. %w", err)
		}

		servers = append(servers, swaggerSrv)
		err = swaggerSrv.Start()

		return fmt.Errorf("swagger server was terminated with an error. %w", err)
	})

	g.Go(func() error {
		restRouter := chi.NewMux()
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

		restAddress := fmt.Sprintf("%v:%v", cfg.Http.Host, cfg.Http.Port)
		restSrv, err := infrastructure.NewServer(logger, restRouter, restAddress, cfg)
		if err != nil {
			return fmt.Errorf("http server creating failed. %w", err)
		}
		servers = append(servers, restSrv)

		err = restSrv.Start()

		return fmt.Errorf("rest server was terminated with an error. %w", err)
	})

	log.Info().
		Str("version", cfg.App.Version).
		Bool("debug", *debug).
		Str("environment", cfg.App.Environment).
		Msgf("Starting services: %s", cfg.App.Name)

	err := g.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("server start failed")
	}
}

func Stop() {
	logger.Warn().Msg("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	for _, server := range servers {
		go func() {
			err := server.Stop(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error while stopping")
			}
		}()
	}

	logger.Warn().Msg("app has stopped")
}
