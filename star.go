package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func getNextPageURL(headerLink string) string {
	var nextPage string
	nextRe := regexp.MustCompile("<([^>]+)>; rel=\"next\"")
	nextMatch := nextRe.FindStringSubmatch(headerLink)
	nextPage = nextMatch[1]

	return nextPage
}

func getLastPageURL(headerLink string) string {
	var lastPage string
	lastRe := regexp.MustCompile("<([^>]+)>; rel=\"last\"")
	lastMatch := lastRe.FindStringSubmatch(headerLink)
	lastPage = lastMatch[1]

	return lastPage
}

const defaultPerPage = 30

func getStarsInfo(repo, token string) (*http.Response, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?per_page=%d&page=503", repo, defaultPerPage)

	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res, nil
}

func RepoStargazers(token string, url string) (*http.Response, error) {
	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res, nil
}

type StarRecord struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func getRepoStarRecords(repo string, token string, maxRequestAmount int) ([]StarRecord, error) {
	starInfo, err := getStarsInfo(repo, token)
	if err != nil {
		return nil, err
	}

	headerLink := starInfo.Header["Link"]
	if headerLink[0] == "" {
		return nil, nil
	}

	for {
		nextPage := getNextPageURL(headerLink[0])
		lastPage := getLastPageURL(headerLink[0])

		fmt.Println(nextPage)
		fmt.Println(lastPage)
		starInfo, err = RepoStargazers(token, nextPage)
		if err != nil {
			return nil, err
		}
		headerLink = starInfo.Header["Link"]
		break
		if headerLink[0] == "" {
			break
		}
	}

	return nil, nil

}
