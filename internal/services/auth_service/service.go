package auth

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authService struct {
	repo interfaces.UserRepo
}

func New(repo interfaces.UserRepo) *authService {
	return &authService{repo: repo}
}

func (as *authService) Authorize(ctx context.Context, uname, pass string) (*models.User, error) {
	return &models.User{
		ID:           primitive.ObjectID{},
		Username:     "test",
		Password:     "",
		Email:        "",
		FirstName:    "",
		LastName:     "",
		CreationDate: 0,
	}, nil
}
