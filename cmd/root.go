package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"star-golang-orms/configs"
	"star-golang-orms/internal"
	"star-golang-orms/pkg"
	"syscall"
)

func Execute() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	// for _, repoNm := range pkg.TargetRepository {
	// 	go internal.GetRepo(repoNm, config.GithubToken)
	// }

	repos, detaiRepos, err := ExecGitHubAPI(config.GithubToken)
	if err != nil {
		log.Println(err)
	}

	err = internal.Edit(repos, detaiRepos)
	if err != nil {
		log.Println(err)
	}
}

func ExecGitHubAPI(token string) ([]internal.GithubRepository, []internal.CheckMouth, error) {
	ctx, cancel := NewCtx()
	defer cancel()

	var repos []internal.GithubRepository
	var detaiRepos []internal.CheckMouth
	for _, repoNm := range pkg.TargetRepository {
		repo, err := internal.NowGithubRepoCount(ctx, repoNm, token)
		if err != nil {
			log.Println(err)
			break
		}
		repos = append(repos, repo)
		detaiRepo, err := internal.GetRepo(ctx, repoNm, token, repo)
		if err != nil {
			log.Println(err)
			break
		}
		detaiRepos = append(detaiRepos, detaiRepo)
	}

	return repos, detaiRepos, nil
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

// 各リポジトリごとにテーブルを作成しし、半年か3ヶ月ごとのスター数の数位をREADMEに

// goroutin を途中キャンセルできるように

// チャート設計

// チャートをREADMEに書き込む
