package user_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type unitTestSuit struct {
	suite.Suite
}

var (
	userName = "test123"
	user     = models.User{
		ID:           primitive.ObjectID{},
		Username:     userName,
		Password:     "$2a$04$3Fwej2KBe58nKVdo0n9mqugGQrEdwzvJqF1JBUgDI3TLLzntYOW96",
		Email:        "user123@ya.ru",
		FirstName:    "test",
		LastName:     "123",
		CreationDate: 0,
	}
	testErr = errors.New("test error")
)

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &unitTestSuit{})
}

func (u *unitTestSuit) TestGetAllEmpty() {
	r := new(repositories.MockUserRepository)

	emptySlice := make([]*models.User, 0)
	r.On("GetAll").Return(emptySlice)

	us := New(r)

	users, err := us.GetAll(context.Background())

	u.Nil(err, "error must be nil")
	u.NotNil(users, "slice must be not nil")
	u.Equal(len(users), 0, fmt.Sprintf("Invalid array length. Except: %v, Receive: %v", 0, len(users)))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestGetAllNotEmpty() {
	r := new(repositories.MockUserRepository)

	userSlice := []*models.User{&user}

	r.On("GetAll").Return(userSlice)

	us := New(r)

	users, err := us.GetAll(context.Background())

	u.Nil(err, "error must be nil")
	u.NotNil(users, "slice must be not nil")
	u.Equal(len(users), 1, fmt.Sprintf("Invalid array length. Except: %v, Receive: %v", 1, len(users)))
	u.Equal(users[0], &user, fmt.Sprintf("Users not equals. Except: %v, Receive: %v", &user, users[0]))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestGetAllError() {
	r := new(repositories.MockUserRepository)
	r.On("GetAll").Return(nil, testErr)

	us := New(r)

	users, err := us.GetAll(context.Background())

	u.NotNil(err, "error must be not nil")
	u.Nil(users, "slice must be nil")
	u.ErrorIs(err, testErr, fmt.Sprintf("error must be testErr. Expected: %s, Received: %s", testErr, err))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestCreateSuccess() {
	r := new(repositories.MockUserRepository)

	r.On("Insert", &user).Return(nil)

	us := New(r)

	err := us.Create(context.Background(), &user)

	u.Nil(err)
}

func (u *unitTestSuit) TestCreateError() {
	r := new(repositories.MockUserRepository)

	r.On("Insert", &user).Return(testErr)

	us := New(r)

	err := us.Create(context.Background(), &user)

	u.NotNil(err, "error must be not nil")
	u.ErrorIs(err, testErr, fmt.Sprintf("error must be testErr. Expected: %s, Received: %s", testErr, err))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestUpdateError() {
	r := new(repositories.MockUserRepository)

	r.On("Update", &user).Return(testErr)

	us := New(r)

	err := us.Update(context.Background(), &user)

	u.NotNil(err, "error must be not nil")
	u.ErrorIs(err, testErr, fmt.Sprintf("error must be testErr. Expected: %s, Received: %s", testErr, err))

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestUpdateSuccess() {
	r := new(repositories.MockUserRepository)

	r.On("Update", &user).Return(nil)

	us := New(r)

	err := us.Update(context.Background(), &user)

	u.Nil(err, "error must be nil")

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestUpdatePasswordSuccess() {
	r := new(repositories.MockUserRepository)

	r.On("UpdatePassword", &user).Return(nil)

	us := New(r)

	err := us.UpdatePassword(context.Background(), &user)

	u.Nil(err, "error must be nil")

	r.AssertExpectations(u.T())
}

func (u *unitTestSuit) TestUpdatePasswordError() {
	r := new(repositories.MockUserRepository)

	r.On("UpdatePassword", &user).Return(testErr)

	us := New(r)

	err := us.UpdatePassword(context.Background(), &user)

	u.NotNil(err, "error must be not nil")
	u.ErrorIs(err, testErr, fmt.Sprintf("error must be testErr. Expected: %s, Received: %s", testErr, err))

	r.AssertExpectations(u.T())
}
