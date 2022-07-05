package auth_service_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/repositories"
	"gitlab.com/g6834/team17/auth-service/internal/services/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

var (
	userName     = "test123"
	userPassword = "qwerty"
	user         = models.User{
		ID:           primitive.ObjectID{},
		Username:     userName,
		Password:     "$2a$04$3Fwej2KBe58nKVdo0n9mqugGQrEdwzvJqF1JBUgDI3TLLzntYOW96",
		Email:        "user123@ya.ru",
		FirstName:    "test",
		LastName:     "123",
		CreationDate: 0,
	}
	jwtSettings = auth_service.JwtSettings{
		SecretKey:  "628f955942efffd7e8e30256",
		AtLifeTime: 5,
		RtLifeTime: 5,
	}
)

type unitTestSuit struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &unitTestSuit{})
}

func (u *unitTestSuit) TestAuthorizeSuccess() {
	r := new(repositories.MockUserRepository)

	r.On("GetByName", userName).Return(&user)

	as := auth_service.New(&jwtSettings, r)

	tokens, err := as.Authorize(context.Background(), userName, userPassword)

	u.Nil(err, "error must be nil")
	u.NotNil(tokens, "tokens must not be nil")
	u.NotNil(tokens.AccessToken, "access token must not be nil")
	u.NotNil(tokens.RefreshToken, "refresh token must not be nil")

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestAuthorizeWrongUsername() {
	r := new(repositories.MockUserRepository)

	r.On("GetByName", userName).Return(nil, errors.New(""))

	as := auth_service.New(&jwtSettings, r)

	tokens, err := as.Authorize(context.Background(), userName, userPassword)

	u.Nil(tokens, "tokens must be nil")
	u.NotNil(err, "error must be not nil")
	u.ErrorIs(err, auth_service.WrongUnameOrPassErr, fmt.Sprintf("error must be WrongUnameOrPassErr. Expected: %s, Received: %s", auth_service.WrongUnameOrPassErr, err))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestAuthorizeWrongPass() {
	r := new(repositories.MockUserRepository)

	r.On("GetByName", userName).Return(&user)

	as := auth_service.New(&jwtSettings, r)

	tokens, err := as.Authorize(context.Background(), userName, userPassword+"x")

	u.Nil(tokens, "tokens must be nil")
	u.NotNil(err, "error must be not nil")
	u.ErrorIs(err, auth_service.WrongUnameOrPassErr, fmt.Sprintf("error must be WrongUnameOrPassErr. Expected: %s, Received: %s", auth_service.WrongUnameOrPassErr, err))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestVerifyTokenSuccess() {
	r := new(repositories.MockUserRepository)

	r.On("GetByName", userName).Return(&user)

	as := auth_service.New(&jwtSettings, r)

	tokens, err := as.Authorize(context.Background(), userName, userPassword)
	u.Nil(err, "error must be nil")

	pair, err := as.VerifyToken(context.Background(), &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})

	u.Nil(err, "error must be nil")
	u.NotNil(pair, "tokens must not be nil")
	u.NotNil(pair.AccessToken, "access token must not be nil")
	u.NotNil(pair.RefreshToken, "refresh token must not be nil")

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestVerifyTokenWrongTokens() {
	r := new(repositories.MockUserRepository)

	r.On("GetByName", userName).Return(&user)

	as := auth_service.New(&jwtSettings, r)

	tokens, err := as.Authorize(context.Background(), userName, userPassword)
	u.Nil(err, "error must be nil")

	t, err := as.VerifyToken(context.Background(), &models.TokenPair{
		AccessToken:  tokens.AccessToken + "x",
		RefreshToken: tokens.RefreshToken,
	})

	u.Nil(t, "access token must be nil")
	u.NotNil(err, "error must be not nil")

	t, err = as.VerifyToken(context.Background(), &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken + "x",
	})

	u.Nil(t, "refresh token must be nil")
	u.NotNil(err, "error must be not nil")
}

func (u *unitTestSuit) TestParseTokenSuccess() {
	r := new(repositories.MockUserRepository)

	r.On("GetByName", userName).Return(&user)

	as := auth_service.New(&jwtSettings, r)

	tokens, err := as.Authorize(context.Background(), userName, userPassword)
	u.Nil(err, "error must be nil")

	us, ok, err := as.ParseToken(context.Background(), tokens.AccessToken)

	u.Nil(err, "error must be nil")
	u.NotNil(us, "user must not be nil")
	u.Equal(ok, true)
	u.Equal(us.Username, user.Username, fmt.Sprintf("wrong user Username. Expected: %s, Received: %s", user.Username, us.Username))
	u.Equal(us.FirstName, user.FirstName, fmt.Sprintf("wrong user FirstName. Expected: %s, Received: %s", user.FirstName, us.FirstName))
	u.Equal(us.LastName, user.LastName, fmt.Sprintf("wrong user LastName. Expected: %s, Received: %s", user.LastName, us.LastName))
	u.Equal(us.Email, user.Email, fmt.Sprintf("wrong user Email. Expected: %s, Received: %s", user.Email, us.Email))
	u.Equal(len(us.Password), 0, "Password must be nil")

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestParseTokenFailed() {
	r := new(repositories.MockUserRepository)

	as := auth_service.New(&jwtSettings, r)

	us, ok, err := as.ParseToken(context.Background(), "something token")

	u.NotNil(err, "error must not be nil")
	u.Nil(us, "user must be nil")
	u.Equal(ok, false)

	r.AssertExpectations(u.T())
}
