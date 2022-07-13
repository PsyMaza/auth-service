package infrastructure

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	db "gitlab.com/g6834/team17/auth-service/internal/databases"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var MigrationNoChange = errors.New("no change")

func NewDatabase(ctx context.Context, cfg *config.Config, log *zerolog.Logger) *mongo.Database {
	mongoCfg := db.MongoConfig{
		Timeout:  cfg.Database.Timeout * int(time.Second),
		DBname:   cfg.Database.Name,
		Username: cfg.Database.User,
		Password: cfg.Database.Password,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
	}

	mongo, err := db.Connect(mongoCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't connection to MongoDB")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err = mongo.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't ping to MongoDB")
	}

	runMigrations(mongo, cfg, log)

	return mongo
}

func runMigrations(mongo *mongo.Database, cfg *config.Config, log *zerolog.Logger) {
	driver, err := mongodb.WithInstance(mongo.Client(), &mongodb.Config{DatabaseName: cfg.Database.Name})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't get databases driver")
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+cfg.Database.Migrations, cfg.Database.Name, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}
	if err := m.Up(); err != nil && errors.Is(err, MigrationNoChange) {
		log.Fatal().Err(err).Msg("Migration failed")
	}
}
