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
	res, err := getStarsForMonthAgo(repo, accessToken)
	// res, err := getRepoLogoUrl(repo, accessToken)
	if err != nil {
		return
	}
	log.Println(res)
}
