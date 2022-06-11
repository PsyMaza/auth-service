package main

import (
	"gitlab.com/g6834/team17/auth-service/internal/config"

	"github.com/rs/zerolog/log"
)

const (
	configName = "config.yaml"
)

func main() {
	if err := config.ReadConfigYML(configName); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.New()

	log.Info().
		Str("version", cfg.Project.Version).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)
	//router.GET("/", func(ctx *gin.Context) {
	//	ctx.JSON(200, gin.H{
	//		"hello": "Hello world !!",
	//	})
	//})
	//log.Fatal(router.Run(":3000"))
}
