package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/g6834/team17/auth_service/internal/db"
	"gitlab.com/g6834/team17/auth_service/internal/handler/api"
	"gitlab.com/g6834/team17/auth_service/internal/repo"
	"gitlab.com/g6834/team17/auth_service/internal/service/auth_service"
	"gitlab.com/g6834/team17/auth_service/internal/service/user_service"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var (
	router = gin.Default()
)

func main() {
	mongo := configureMongo()
	userRepo := repo.New(mongo)
	userService := user_service.New(userRepo)
	authService := auth_service.New("628f955942efffd7e8e30256", userRepo)

	userCrud := api.NewUserCrudHandler(userService)
	auth := api.NewAuthorizeHandler(authService)

	router.GET("/getAll", userCrud.GetAll)
	router.POST("/login", auth.Login)
	router.POST("/create", userCrud.Create)
	router.POST("/update", userCrud.Update)
	log.Fatal(router.Run(":8080"))
}

func configureMongo() *mongo.Database {
	mongoCfg := db.MongoConfig{
		Timeout:  5000,
		DBname:   "mts",
		Username: "",
		Password: "",
		Host:     "0.0.0.0",
		Port:     "27017",
	}

	mongo, err := db.Connect(mongoCfg)
	if err != nil {
		log.Fatal(err)
	}

	return mongo
}
