package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pr-manager-service/config"
	"pr-manager-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err = app.Run(ctx, cfg)
	if err != nil {
		log.Fatalf("Application run error: %v", err)
	}
}
