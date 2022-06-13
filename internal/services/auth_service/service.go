package auth_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type JwtSettings struct {
	SecretKey  string
	AtLifeTime int
	RtLifeTime int
}

type authService struct {
	jwtSettings *JwtSettings
	repo        interfaces.UserRepo
}

var NotFoundUserErr = errors.New("no user found with this username and password")

const (
	authorized = "authorized"
	userId     = "user_id"
	expired    = "expired"
	email      = "email"
	firstName  = "first_name"
	lastName   = "last_name"
	username   = "username"
)

func New(jwtSettings *JwtSettings, repo interfaces.UserRepo) *authService {
	return &authService{repo: repo, jwtSettings: jwtSettings}
}

func (as *authService) Authorize(ctx context.Context, uname, pass string) (*models.TokenDetails, error) {
	user, err := as.repo.GetByName(ctx, uname)
	if err != nil {
		log.Println(err)
		return nil, NotFoundUserErr
	}

	err = utils.CheckPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		return nil, err
	}

	token, err := createToken(user, as.jwtSettings)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func createToken(user *models.User, settings *JwtSettings) (td *models.TokenDetails, err error) {
	td = &models.TokenDetails{
		AtExpires: time.Now().Add(time.Minute * time.Duration(settings.AtLifeTime)),
		RtExpires: time.Now().Add(time.Hour * time.Duration(settings.RtLifeTime)),
	}

	atClaims := jwt.MapClaims{}
	atClaims[authorized] = true
	atClaims[userId] = user.ID
	atClaims[username] = user.Username
	atClaims[firstName] = user.FirstName
	atClaims[email] = user.Email
	atClaims[lastName] = user.LastName
	atClaims[expired] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(settings.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("get access token error: %w", err)
	}

	rtClaims := jwt.MapClaims{}
	rtClaims[userId] = user.ID
	rtClaims[expired] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(settings.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("get refresh token error: %w", err)
	}

	return
}

func (as *authService) VerifyToken(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.jwtSettings.SecretKey), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
func (as *authService) ParseToken(ctx context.Context, tokenString string) (*models.User, bool, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.jwtSettings.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, false, err
	}

	stringId := fmt.Sprintf("%v", claims[userId])
	id, _ := primitive.ObjectIDFromHex(stringId)
	user := &models.User{
		ID:        id,
		Username:  fmt.Sprintf("%v", claims[username]),
		Email:     fmt.Sprintf("%v", claims[email]),
		FirstName: fmt.Sprintf("%v", claims[firstName]),
		LastName:  fmt.Sprintf("%v", claims[lastName]),
	}

	return user, true, nil
}
