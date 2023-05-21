package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Stargazer struct {
	StarredAt time.Time `json:"starred_at"`
}

type GithubRepository struct {
	FullName         string `json:"full_name"`
	StargazersCount  int    `json:"stargazers_count"`
	CreatedAt        string `json:"created_at"`
	SubscribersCount int    `json:"subscribers_count"`
	ForksCount       int    `json:"forks_count"`
}

var (
	errNoMorePages  = errors.New("no more pages to get")
	ErrTooManyStars = errors.New("repo has too many stargazers, github won't allow us to list all stars")
)

func getRepo(name, token string) (*http.Response, error) {
	ctx, cancel := NewCtx()
	defer cancel()

	repo, err := nowGithubRepoCount(name, token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(repo) //{uptrace/bun 2076 2021-05-03T11:40:52Z 24 157}

	sem := make(chan bool, 4)
	var eg errgroup.Group
	var lock sync.Mutex
	var stargazers []Stargazer
	for page := 1; page <= lastPage(*repo); page++ { //298回非同期する
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result, err := getStargazersPage(ctx, *repo, page, token)
			if errors.Is(err, errNoMorePages) {
				log.Println(err)
				return nil
			}
			if err != nil {
				log.Println(err)
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			stargazers = append(stargazers, result...)
			return nil
		})
	}

	// stargazers = append(stargazers, [{2017-07-07 02:50:15 +0000 UTC}])
	log.Println(stargazers) //[{2017-07-07 02:50:15 +0000 UTC} {2017-07-07 05:06:33 +0000 UTC} {2017-07-07 10:56:49 +0000 UTC} {2017-07-07 11:25:36 +0000 UTC} {2017-07-07 19:42:38 +0000 UTC} {2017-07-08 01:06:01 +0000 UTC} ]
	return nil, nil
}

func nowGithubRepoCount(repoName, token string) (*GithubRepository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var repo GithubRepository
	if res.StatusCode == http.StatusOK {
		if err := json.Unmarshal(bts, &repo); err != nil {
			return nil, err
		}
	}

	return &repo, nil
}

func stringToTime(str string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", str)
	return t
}

func getStargazersPage(ctx context.Context, repo GithubRepository, page int, token string) ([]Stargazer, error) {
	var stars []Stargazer

	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?per_page=100&page=%d&", repo.FullName, page)
	client := NewHttpClient(url, http.MethodGet, token)
	resp, err := client.SendRequest()
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return stars, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(bts, &stars); err != nil {
			return nil, err
		}
		if len(stars) == 0 {
			return nil, fmt.Errorf("スターなし")
		}
		log.Println("-----")
		log.Println(stars)
		log.Println("-----")
		return stars, nil
	default:
		return nil, fmt.Errorf("その他のエラー")
	}
}

type GithubUser struct {
	AvatarURL string `json:"avatar_url"`
}

func getRepoLogoUrl(repoName string, token string) (string, error) {
	owner := strings.Split(repoName, "/")[0]
	url := fmt.Sprintf("https://api.github.com/users/%s", owner)
	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return "", err
	}

	var user GithubUser
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return "", err
	}
	defer res.Body.Close()

	return user.AvatarURL, nil
}

func totalPages(repo GithubRepository) int {
	pageSize := 100
	return repo.StargazersCount / pageSize
}

func lastPage(repo GithubRepository) int {
	return totalPages(repo) + 1
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
