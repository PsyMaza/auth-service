package test

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	dbName        = "mts"
	migrationFile = "migrations"
)

func TestContainers(t *testing.T) {
	suite.Run(t, new(TestContainersSuite))
}

type TestContainersSuite struct {
	suite.Suite

	//server         *http.Server
	userRepo       interfaces.UserRepo
	mongoContainer testcontainers.Container
}

func (suite *TestContainersSuite) SetupSuite() {
	ctx := context.Background()

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
			SkipReaper:   true,
			AutoRemove:   true,
		},
		Started: true,
	})
	suite.Require().NoError(err)

	// with a second delay migrations work properly
	time.Sleep(time.Second * 5)

	ip, err := dbContainer.Host(ctx)
	suite.Require().NoError(err)
	port, err := dbContainer.MappedPort(ctx, "27017")
	suite.Require().NoError(err)

	clientUrl := fmt.Sprintf("mongodb://%v:%v",
		ip,
		uint16(port.Int()),
	)

	clientOptions := options.Client().ApplyURI(clientUrl)

	client, err := mongo.NewClient(clientOptions)
	suite.Require().NoError(err)

	err = client.Connect(ctx)
	suite.Require().NoError(err)

	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: dbName})
	suite.Require().NoError(err)
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationFile, dbName, driver)
	suite.Require().NoError(err)
	err = m.Up()
	suite.Require().NoError(err)

	userRepo := repositories.NewDatabaseRepo(client.Database(dbName))

	suite.userRepo = userRepo
	suite.mongoContainer = dbContainer

	suite.T().Log("Suite setup is done")
}

func (s *TestContainersSuite) TearDownSuite() {
	s.mongoContainer.Terminate(context.Background())
	s.T().Log("Suite stop is done")
}
