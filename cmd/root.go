package cmd

import (
	"context"
	"log"
	"os/signal"
	"star-golang-orms/app"
	"star-golang-orms/domain/service"
	"star-golang-orms/infra"
	"syscall"
)

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := infra.Load(ctx); err != nil {
		log.Fatal("cannot load config", err)
	}

	svc := setupJob(ctx)
	if err := svc.Start(ctx); err != nil {
		log.Fatal("cannot start job", err)
	}
}

func setupJob(ctx context.Context) service.Fetcher {
	return app.NewFetchService(
		infra.NewGitHubRepository(ctx),
	)
}
