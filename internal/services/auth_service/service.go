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

var WrongUnameOrPassErr = errors.New("no user found with this username and password")

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
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	user, err := as.repo.GetByName(ctx, uname)
	if err != nil {
		log.Println(err)
		return nil, WrongUnameOrPassErr
	}

	err = utils.CheckPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		log.Println(err)
		return nil, WrongUnameOrPassErr
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

func (as *authService) VerifyToken(ctx context.Context, tokens *models.TokenPair) (*models.TokenPair, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	at, err := parseToken(tokens.AccessToken, as.jwtSettings.SecretKey)
	if err != nil {
		return nil, err
	}

	if !at.Valid {
		return as.newPairToken(ctx, tokens.AccessToken)
	}

	rt, err := parseToken(tokens.RefreshToken, as.jwtSettings.SecretKey)
	if err != nil {
		return nil, err
	}

	if !rt.Valid {
		return as.newPairToken(ctx, tokens.RefreshToken)
	}

	return &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (as *authService) newPairToken(ctx context.Context, token string) (*models.TokenPair, error) {
	user, _, err := as.ParseToken(ctx, token)
	if err != nil {
		return nil, err
	}

	td, err := createToken(user, as.jwtSettings)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}, nil
}

func parseToken(token, secretKey string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (as *authService) ParseToken(ctx context.Context, tokenString string) (*models.User, bool, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

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
