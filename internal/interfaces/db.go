package interfaces

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/models"
)

type UserRepo interface {
	Get(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByName(ctx context.Context, uname string) (*models.User, error)
	Insert(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, user *models.User) error
}
