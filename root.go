package main

import (
	"log"
)

func Execute() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	// repoName := "beego/beego"
	repoName := "uptrace/bun"
	// repoName := ""
	accessToken := config.GithubToken

	// res, err := getRepoStargazers(repo, accessToken, 1)
	// res, err := getRepoStargazersCount(repo, accessToken)
	// res, err := getStarsInfo(repo, accessToken)
	// res, err := getRepoStarRecords(repo, accessToken, 1)
	res, err := getRepoLogoUrl(repoName, accessToken)
	// res, err := nowGithubRepoCount(repoName, accessToken)
	// res, err := getRepo(name, accessToken)
	if err != nil {
		return
	}
	log.Println(res)
}

//正確な値がとれないなら1ヶ月に1回取り込んで、READMEに書き込む？
//最終的に値がほしいリポジトリをまとめてそれをgoroutinでとってくる
