package user_service

import (
	"context"
	"gitlab.com/g6834/team17/auth_service/internal/handler/api"
	"gitlab.com/g6834/team17/auth_service/internal/helper"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo repo.UserRepo
}

func New(repo repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetByNameAndPass(ctx context.Context, uname, pass string) (*api.User, error) {
	user, err := us.repo.GetByNameAndPass(ctx, uname, pass)
	if err != nil {
		return nil, err
	}
	return toApiUser(user), nil
}

func (us *UserService) Create(ctx context.Context, user *api.User) (err error) {
	err = us.repo.Insert(ctx, toDbUser(user))
	return err
}

func (us *UserService) Update(ctx context.Context, user *api.User) (err error) {
	err = us.repo.Insert(ctx, toDbUser(user))
	return err
}

func toApiUser(user *model.User) *api.User {
	return &api.User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func toDbUser(user *api.User) *model.User {
	id, _ := primitive.ObjectIDFromHex(user.ID)

	return &model.User{
		ID:        id,
		Username:  user.Username,
		Email:     user.Email,
		Password:  helper.GetHash([]byte(user.Password)),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
