package user_service

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/api/handlers"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	repo repo.UserRepo
}

func New(repo repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Create(user *handlers.User) (err error) {
	err = us.repo.Insert(context.Background(), toDbUser(user))
	return err
}

func (us *UserService) Update(user *handlers.User) (err error) {
	err = us.repo.Insert(context.Background(), toDbUser(user))
	return err
}

func toDbUser(user *handlers.User) *models.User {
	id, _ := primitive.ObjectIDFromHex(user.ID)

	return &models.User{
		ID:        id,
		Username:  user.Username,
		Email:     user.Email,
		Password:  getHash([]byte(user.Password)),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
