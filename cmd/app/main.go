package main

import (
	"context"
	"os/signal"
	"syscall"

	"url-shortner/internal/app"
	"url-shortner/internal/app/config"
)

func main() {
	cfg := config.NewConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx, cfg)
}
