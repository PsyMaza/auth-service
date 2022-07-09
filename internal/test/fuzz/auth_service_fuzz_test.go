package fuzz

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/repositories"
	"gitlab.com/g6834/team17/auth-service/internal/services/user_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	dbName        = "mts"
	migrationFile = "migrations"
)

type dbContainer struct {
	testcontainers.Container
	mongo *mongo.Client
}

func setupDbContainer() (*dbContainer, error) {
	ctx := context.Background()

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
			SkipReaper:   true,
			AutoRemove:   true,
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := mongoContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	clientUrl := fmt.Sprintf("mongodb://%v:%v",
		ip,
		uint16(port.Int()),
	)

	clientOptions := options.Client().ApplyURI(clientUrl)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: dbName})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationFile, dbName, driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil {
		return nil, err
	}

	return &dbContainer{
		Container: mongoContainer,
		mongo:     client,
	}, nil
}

func FuzzTestCreateUser(f *testing.F) {
	ctx := context.Background()

	db, err := setupDbContainer()
	if err != nil {
		f.Error(err)
	}
	defer func() {
		err := db.Container.Terminate(ctx)
		if err != nil {
			f.Log(err)
		}
	}()

	userR := repositories.NewDatabaseRepo(db.mongo.Database(dbName))
	userS := user_service.New(userR)

	f.Fuzz(func(t *testing.T, username, pass, email, firstName, lastName string) {
		user := &models.User{
			ID:           primitive.ObjectID{},
			Username:     username,
			Password:     pass,
			Email:        email + "ya.ru",
			FirstName:    firstName,
			LastName:     lastName,
			CreationDate: uint64(time.Now().Unix()),
		}

		err := userS.Create(ctx, user)

		if err != nil && !mongo.IsDuplicateKeyError(err) {
			t.Error(err)
		}
	})
}
