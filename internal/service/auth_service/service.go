package auth_service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/g6834/team17/auth_service/internal/handler/api"
	"gitlab.com/g6834/team17/auth_service/internal/service/user_service"
	"log"
	"time"
)

type AuthService struct {
	SecretKey string
	us        user_service.UserService
}

var NotFoundUserErr = errors.New("No user found with this username and password")

func New(secretKey string, us user_service.UserService) *AuthService {
	return &AuthService{
		secretKey,
		us,
	}
}

var (
	authorized = "authorized"
	userId     = "user_id"
	exp        = "exp"
)

func (as *AuthService) Authorize(ctx context.Context, uname, pass string) (string, error) {
	user, err := as.us.GetByNameAndPass(ctx, uname, pass)
	if err != nil {
		log.Println(err)
		return "", NotFoundUserErr
	}

	token, err := createToken(user, as.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createToken(user *api.User, secretKey string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims[authorized] = true
	atClaims[userId] = user.ID
	atClaims[exp] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
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
