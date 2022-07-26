package infrastructure

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	l      net.Listener
	port   int
}

type ServerConfig struct {
}

func NewServer(log *zerolog.Logger, handler http.Handler, address string, cfg *config.Config) (*Server, error) {
	var (
		err error
		s   Server
	)
	s.l, err = net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed listen port")
	}
	s.port = s.l.Addr().(*net.TCPAddr).Port

	s.server = &http.Server{
		Handler:      handler,
		ReadTimeout:  time.Second * time.Duration(cfg.Http.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Http.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Http.IdleTimeout),
	}

	return &s, nil
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.l); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
