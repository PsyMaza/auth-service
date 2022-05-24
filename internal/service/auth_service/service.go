package auth_service

import (
	"errors"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"time"
)

type AuthService struct {
	SecretKey string
}

func New(secretKey string) *AuthService {
	return &AuthService{
		secretKey,
	}
}

var (
	authorized = "authorized"
	userId     = "user_id"
	exp        = "exp"
)

func (s *AuthService) Authorize(uname, pass string) (string, error) {
	if uname != validUser.Username && pass != validUser.Password {
		return "", errors.New("Invalid login or password")
	}

	return s.createToken(&validUser)
}

func (s *AuthService) createToken(user *model.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims[authorized] = true
	atClaims[userId] = user.ID
	atClaims[exp] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	token, err := at.SignedString([]byte(s.SecretKey))
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

var validUser = model.User{
	ID:       1,
	Username: "test",
	Password: "123",
}
