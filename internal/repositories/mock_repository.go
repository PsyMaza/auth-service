package repositories

import (
	"context"
	"github.com/stretchr/testify/mock"
	"gitlab.com/g6834/team17/auth-service/internal/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	args := r.Called()

	return args.Get(0).([]*models.User), nil
}

func (r *MockUserRepository) Get(ctx context.Context, id string) (*models.User, error) {
	args := r.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), nil
}

func (r *MockUserRepository) GetByName(ctx context.Context, uname string) (*models.User, error) {
	args := r.Called(uname)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), nil
}

func (r *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	r.Called(user)

	return nil
}

func (r *MockUserRepository) Insert(ctx context.Context, user *models.User) error {
	r.Called(user)

	return nil
}

func (r *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	r.Called(user)

	return nil
}

func (r *MockUserRepository) UpdatePassword(ctx context.Context, user *models.User) error {
	r.Called(user)

	return nil
}
