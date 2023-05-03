package main

import (
	"log"
)

func Execute() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	repo := "beego/beego"
	accessToken := config.GithubToken

	// res, err := getRepoStargazers(repo, accessToken, 1)
	// res, err := getRepoStargazersCount(repo, accessToken)
	// res, err := getStarsInfo(repo, accessToken)
	res, err := getRepoStarRecords(repo, accessToken, 1)
	// res, err := getRepoLogoUrl(repo, accessToken)
	if err != nil {
		return
	}
	log.Println(res)
}
