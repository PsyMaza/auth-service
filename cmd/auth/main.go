package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/g6834/team17/auth_service/internal/handler/api"
	"log"
)

var (
	router = gin.Default()
)

func main() {
	router.POST("/login", api.Login)
	log.Fatal(router.Run(":8080"))
}
