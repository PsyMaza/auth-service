package main

import (
	"gitlab.com/g6834/team17/auth-service/config"
	"log"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()	
	log.Println(cfg)
	//router.GET("/", func(ctx *gin.Context) {
	//	ctx.JSON(200, gin.H{
	//		"hello": "Hello world !!",
	//	})
	//})
	//log.Fatal(router.Run(":3000"))
}
