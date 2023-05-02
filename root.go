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

	res, err := GetRepoStargazers(repo, accessToken, 1)
	// res, err := GetRepoStargazersCount(repo, accessToken)
	log.Println(res)
	if err != nil {
		return
	}
}
