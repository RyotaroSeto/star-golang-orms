package cmd

import (
	"log"
	"star-golang-orms/configs"
	"star-golang-orms/internal"
)

func Execute() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	// repoName := "beego/beego"
	repoName := "uptrace/bun"
	// repoName := ""
	accessToken := config.GithubToken

	// res, err := getRepoStarRecords(repo, accessToken, 1)
	// res, err := getRepoLogoUrl(repoName, accessToken)
	// res, err := nowGithubRepoCount(repoName, accessToken)
	// res, err := getRepo(repoName, accessToken)

	ctx, cancel := internal.NewCtx()
	defer cancel()
	repo, err := internal.NowGithubRepoCount(repoName, accessToken)
	if err != nil {
		log.Println(err)
	}
	res, err := internal.GetStargazersPage(ctx, *repo, 10, accessToken)
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
