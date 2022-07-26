package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.com/g6834/team17/auth-service/internal/application"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM|syscall.SIGINT, os.Interrupt)
	defer cancel()

	go application.Start(ctx)
	<-ctx.Done()
	application.Stop()
}
