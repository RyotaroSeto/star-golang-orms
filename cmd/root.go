package cmd

import (
	"log"
	"star-golang-orms/configs"
	"star-golang-orms/internal"
	"star-golang-orms/pkg"
)

func Execute() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	targetRepositorys := pkg.TargetRepository

	// res, err := getRepoStarRecords(repo, config.GithubToken, 1)
	// res, err := getRepoLogoUrl(targetRepositorys[0], config.GithubToken)
	// res, err := nowGithubRepoCount(targetRepositorys[0], config.GithubToken)
	// res, err := getRepo(targetRepositorys[0], config.GithubToken)

	ctx, cancel := internal.NewCtx()
	defer cancel()
	repo, err := internal.NowGithubRepoCount(targetRepositorys[0], config.GithubToken)
	if err != nil {
		log.Println(err)
	}
	res, err := internal.GetStargazersPage(ctx, *repo, 10, config.GithubToken)
	if err != nil {
		return
	}
	log.Println(res)
}

//正確な値がとれないなら1ヶ月に1回取り込んで、READMEに書き込む？
//最終的に値がほしいリポジトリをまとめてそれをgoroutinでとってくる

// 1リポジトリ以下をREADMEに書き込む
// FullName         string `json:"full_name"`
// StargazersCount  int    `json:"stargazers_count"`
// CreatedAt        string `json:"created_at"`

// 複数リポジトリREADMEに書き込む

// goroutin を途中キャンセルできるように

// チャート設計

// チャートをREADMEに書き込む

// 現状、スター数のみだが他のカウントも取得するようにする
// 1リポジトリ以下をREADMEに書き込む
// FullName         string `json:"full_name"`
// StargazersCount  int    `json:"stargazers_count"`
// CreatedAt        string `json:"created_at"`
// SubscribersCount int    `json:"subscribers_count"`
// ForksCount       int    `json:"forks_count"`
