package presenters

import (
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
)

type presenters struct {
	logger zerolog.Logger
}

func NewPresenters(logger zerolog.Logger) interfaces.Presenters {
	return &presenters{logger: logger}
}
