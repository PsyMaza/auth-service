package auth_service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/g6834/team17/auth_service/internal/helper"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/service/user_service"
	"log"
	"time"
)

type AuthService interface {
	Authorize(ctx context.Context, uname, pass string) (string, error)
}

type authService struct {
	SecretKey string
	us        user_service.UserService
}

var NotFoundUserErr = errors.New("No user found with this username and password")

func New(secretKey string, us user_service.UserService) AuthService {
	return &authService{
		secretKey,
		us,
	}
}

var (
	authorized = "authorized"
	userId     = "user_id"
	exp        = "exp"
	email      = "email"
	firstName  = "first_name"
	lastName   = "last_name"
	username   = "username"
)

func (as *authService) Authorize(ctx context.Context, uname, pass string) (string, error) {
	user, err := as.us.GetByName(ctx, uname)
	if err != nil {
		log.Println(err)
		return "", NotFoundUserErr
	}

	err = helper.CheckPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		return "", err
	}

	token, err := createToken(user, as.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createToken(user *model.User, secretKey string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims[authorized] = true
	atClaims[userId] = user.ID
	atClaims[username] = user.Username
	atClaims[firstName] = user.FirstName
	atClaims[lastName] = user.LastName
	atClaims[exp] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

//func (s *AuthService) Validate(token string) (bool, error) {
//
//}
//
//func (s *AuthService) Unathorize(token string) error {
//
//}
//func (s *AuthService) Refresh(accessToken, refreshToken string) error {
//
//}
