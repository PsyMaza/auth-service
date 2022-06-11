package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/config"
	"gitlab.com/g6834/team17/auth-service/internal/database"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/repo"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	configName = "config.yaml"
)

var MigrationNoChange = errors.New("no change")

func main() {
	if err := config.ReadConfigYML(configName); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.New()

	debug := flag.Bool("debug", cfg.App.Debug, "sets log level to debug")
	useDatabase := flag.Bool("use-db", true, "sets repository source")
	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

	log.Info().
		Str("version", cfg.App.Version).
		Bool("debug", *debug).
		Str("environment", cfg.App.Environment).
		Msgf("Starting service: %s", cfg.App.Name)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var userRepo interfaces.UserRepo

	mongo := configureMongo(&cfg)
	*useDatabase = false // todo: delete when will be created mongodb in k8s
	if *useDatabase {
		userRepo = repo.NewDatabaseRepo(mongo)
	} else {
		userRepo = nil
	}

	fmt.Println(userRepo)

	//*migration = false // todo: delete when will be created mongodb in k8s
	if *migration {
		runMigrations(mongo, cfg)
	}
}

func runMigrations(mongo *mongo.Database, cfg config.Config) {
	driver, err := mongodb.WithInstance(mongo.Client(), &mongodb.Config{DatabaseName: cfg.Database.Name})
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't get database driver")
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+cfg.Database.Migrations, cfg.Database.Name, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}
	if err := m.Up(); err != nil && errors.Is(err, MigrationNoChange) {
		log.Fatal().Err(err).Msg("Migration failed")
	}
}

func configureMongo(cfg *config.Config) *mongo.Database {
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

	return mongo
}
