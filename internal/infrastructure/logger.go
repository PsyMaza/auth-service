package infrastructure

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"os"
	"time"
)

func NewLogger(cfg *config.Config, debug bool) *zerolog.Logger {
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
		Str("host_ip", cfg.Http.Host).
		Str("host_port", fmt.Sprint(cfg.Http.Port)).
		Logger()

	return &logger
}
