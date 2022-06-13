package user_service

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/utils"
)

type userService struct {
	repo interfaces.UserRepo
}

func New(repo interfaces.UserRepo) *userService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := us.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *userService) Create(ctx context.Context, user *models.User) (err error) {
	user.Password = utils.GetHash([]byte(user.Password))
	err = us.repo.Insert(ctx, user)
	return err
}

func (us *userService) Update(ctx context.Context, user *models.User) (err error) {
	err = us.repo.Update(ctx, user)
	return err
}

func (us *userService) UpdatePassword(ctx context.Context, user *models.User) (err error) {
	err = us.repo.Update(ctx, user)
	return err
}
