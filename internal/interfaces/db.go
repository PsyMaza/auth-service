package interfaces

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/model"
)

type UserRepo interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}
