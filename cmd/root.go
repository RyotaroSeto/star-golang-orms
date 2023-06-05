package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"star-golang-orms/configs"
	"star-golang-orms/pkg"
	"syscall"
)

func Execute() {
	ctx, cancel := NewCtx()
	defer cancel()

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
		return
	}

	gh, err := pkg.ExecGitHubAPI(ctx, config.GithubToken)
	if err != nil {
		log.Fatal("cannot exec github api", err)
		return
	}

	err = gh.SortDesByStarCount()
	if err != nil {
		log.Fatal("cannot sort star count", err)
		return
	}

	err = gh.MakeChart()
	if err != nil {
		log.Fatal("connot make chart", err)
		return
	}

	err = gh.Edit()
	if err != nil {
		log.Fatal("connot edit readme", err)
		return
	}
}

func NewCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
		<-trap
	}()

	return ctx, cancel
}
