package interfaces

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/models"
)

type AuthService interface {
	Authorize(ctx context.Context, uname, pass string) (*models.TokenDetails, error)
}

type UserService interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) (err error)
	Update(ctx context.Context, user *models.User) (err error)
	UpdatePassword(ctx context.Context, user *models.User) (err error)
}
