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
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	gh, err := ExecGitHubAPI(config.GithubToken)
	if err != nil {
		log.Println(err)
	}

	err = gh.Edit()
	if err != nil {
		log.Println(err)
	}
}

func ExecGitHubAPI(token string) (pkg.GitHub, error) {
	ctx, cancel := NewCtx()
	defer cancel()

	var repos []pkg.GithubRepository
	var detaiRepos []pkg.ReadmeDetailsRepository
	for _, repoNm := range pkg.TargetRepository {
		log.Println("start:" + repoNm)
		repo, err := pkg.NowGithubRepoCount(ctx, repoNm, token)
		if err != nil {
			log.Println(err)
			break
		}
		repos = append(repos, repo)
		detaiRepo, err := pkg.GetRepo(ctx, repoNm, token, repo)
		if err != nil {
			log.Println(err)
			break
		}
		detaiRepos = append(detaiRepos, detaiRepo)
	}

	gh := pkg.NewGitHub(repos, detaiRepos)
	return gh, nil
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

// リポジトリ作成日より前だったら「-」と出力する

// 取得したリポジトリをスター数順にソート

// 詳細テーブルはリポジトリ作成日からスタートにする

// goroutin を途中キャンセルできるように

// チャート設計

// チャートをREADMEに書き込む
