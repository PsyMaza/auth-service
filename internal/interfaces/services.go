package interfaces

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/models"
)

type AuthService interface {
	Authorize(ctx context.Context, uname, pass string) (*models.User, error)
}
