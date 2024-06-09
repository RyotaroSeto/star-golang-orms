package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"star-golang-orms/app"
	"star-golang-orms/domain/service"
	"star-golang-orms/infra"
	"syscall"
)

func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := infra.Load(ctx); err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	svc := setupJob(ctx)
	if err := svc.Start(ctx); err != nil {
		return fmt.Errorf("cannot start job: %w", err)
	}

	return nil
}

func setupJob(ctx context.Context) service.Fetcher {
	return app.NewFetchService(
		infra.NewGitHubRepository(ctx),
	)
}
