package user_service

import (
	"context"
	"gitlab.com/g6834/team17/auth_service/internal/helper"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/repo"
)

type UserService interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	Create(ctx context.Context, user *model.User) (err error)
	Update(ctx context.Context, user *model.User) (err error)
	UpdatePassword(ctx context.Context, user *model.User) (err error)
}

type userService struct {
	repo repo.UserRepo
}

func New(repo repo.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAll(ctx context.Context) ([]*model.User, error) {
	users, err := us.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *userService) Create(ctx context.Context, user *model.User) (err error) {
	user.Password = helper.GetHash([]byte(user.Password))
	err = us.repo.Insert(ctx, user)
	return err
}

func (us *userService) Update(ctx context.Context, user *model.User) (err error) {
	err = us.repo.Update(ctx, user)
	return err
}

func (us *userService) UpdatePassword(ctx context.Context, user *model.User) (err error) {
	err = us.repo.Update(ctx, user)
	return err
}
