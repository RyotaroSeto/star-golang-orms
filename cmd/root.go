package cmd

import (
	"context"
	"log"
	"os/signal"
	"star-golang-orms/infra"
	"syscall"
)

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := infra.Load(".")
	if err != nil {
		log.Fatal("cannot load config", err)
		return
	}

	log.Println(ctx)
	// gh, err := pkg.ExecGitHubAPI(ctx, config.GitHubToken)
	// if err != nil {
	// 	log.Fatal("cannot exec github api", err)
	// 	return
	// }

	// err = gh.SortDesByStarCount()
	// if err != nil {
	// 	log.Fatal("cannot sort star count", err)
	// 	return
	// }

	// err = gh.MakeChart()
	// if err != nil {
	// 	log.Fatal("cannot make chart", err)
	// 	return
	// }

	// err = gh.Edit()
	// if err != nil {
	// 	log.Fatal("cannot edit readme", err)
	// 	return
	// }
}
