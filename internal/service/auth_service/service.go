package auth_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/g6834/team17/auth_service/internal/helper"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/repo"
	"log"
	"time"
)

type AuthService interface {
	Authorize(ctx context.Context, uname, pass string) (*model.TokenDetails, error)
}

type authService struct {
	SecretKey string
	repo      repo.UserRepo
}

var NotFoundUserErr = errors.New("No user found with this username and password")

func New(secretKey string, repo repo.UserRepo) AuthService {
	return &authService{
		secretKey,
		repo,
	}
}

var (
	authorized = "authorized"
	userId     = "user_id"
	expired    = "expired"
	email      = "email"
	firstName  = "first_name"
	lastName   = "last_name"
	username   = "username"
)

func (as *authService) Authorize(ctx context.Context, uname, pass string) (*model.TokenDetails, error) {
	user, err := as.repo.GetByName(ctx, uname)
	if err != nil {
		log.Println(err)
		return nil, NotFoundUserErr
	}

	err = helper.CheckPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		return nil, err
	}

	token, err := createToken(user, as.SecretKey)
	if err != nil {
		return nil, err
	}

	as.repo.SetToken(ctx, user.ID, token)

	return token, nil
}

func createToken(user *model.User, secretKey string) (td *model.TokenDetails, err error) {
	td = &model.TokenDetails{
		AtExpires: time.Now().Add(time.Minute * 15).Unix(),
		RtExpires: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	atClaims := jwt.MapClaims{}
	atClaims[authorized] = true
	atClaims[userId] = user.ID
	atClaims[username] = user.Username
	atClaims[firstName] = user.FirstName
	atClaims[lastName] = user.LastName
	atClaims[expired] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(secretKey))

	rtClaims := jwt.MapClaims{}
	rtClaims[userId] = user.ID
	rtClaims[expired] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return
}

func (as *authService) VerifyToken(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.SecretKey), nil
	})
	if err != nil {
		return false, err
	}

	err = as.repo.TokenExist(ctx, tokenString, repo.Refresh)
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

//
//func (s *AuthService) Unathorize(token string) error {
//
//}
//func (s *AuthService) Refresh(accessToken, refreshToken string) error {
//
//}
