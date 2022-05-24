package user_service

import (
	"context"
	"gitlab.com/g6834/team17/auth_service/internal/handler/api"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/repo"
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

func (us *UserService) Create(user *api.User) (err error) {
	err = us.repo.Insert(context.Background(), toDbUser(user))
	return err
}

func (us *UserService) Update(user *api.User) (err error) {
	err = us.repo.Insert(context.Background(), toDbUser(user))
	return err
}

func toDbUser(user *api.User) *model.User {
	id, _ := primitive.ObjectIDFromHex(user.ID)

	return &model.User{
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
