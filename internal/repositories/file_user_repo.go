package repositories

import (
	"context"
	"errors"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type FileRepo struct {
	users []*models.User `yaml:"users"`
}

const (
	FILE_PATH = "users.yaml"
)

func NewFileRepo() *FileRepo {
	file, err := os.Open(filepath.Clean(FILE_PATH))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	repo := &FileRepo{}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&repo.users); err != nil {
		panic(err)
	}

	return repo
}

var NotFoundUserErr = errors.New("no user found with this username and password")

func (fr *FileRepo) Get(ctx context.Context, userId string) (*models.User, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	for _, user := range fr.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, NotFoundUserErr
}

func (fr *FileRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	return fr.users, nil
}

func (fr *FileRepo) GetByName(ctx context.Context, uname string) (*models.User, error) {
	for _, user := range fr.users {
		if user.Username == uname {
			return user, nil
		}
	}

	return nil, NotFoundUserErr
}

func (fr *FileRepo) Insert(ctx context.Context, user *models.User) error {
	panic("Not emplement")
}

func (fr *FileRepo) Update(ctx context.Context, user *models.User) error {
	panic("Not emplement")
}

func (fr *FileRepo) UpdatePassword(ctx context.Context, user *models.User) error {
	panic("Not emplement")
}
