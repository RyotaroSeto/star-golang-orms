package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	baseURL   = "https://api.github.com/"
	rateLimit = "rate_limit"
)

type Stargazer struct {
	StarredAt time.Time `json:"starred_at"`
}

type GithubRepository struct {
	FullName         string    `json:"full_name"`
	URL              string    `json:"html_url"`
	Description      string    `json:"description"`
	StargazersCount  int       `json:"stargazers_count"`
	SubscribersCount int       `json:"subscribers_count"`
	ForksCount       int       `json:"forks_count"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ReadmeDetailsRepository struct {
	RepoName   string
	RepoURL    string
	StarCounts map[string]int
}

func ExecGitHubAPI(ctx context.Context, token string) (GitHub, error) {
	var repos []GithubRepository
	var detaiRepos []ReadmeDetailsRepository

	wg := new(sync.WaitGroup)
	var lock sync.Mutex
	for _, repoNm := range TargetRepository {
		wg.Add(1)
		go func(repoNm string) {
			defer wg.Done()
			repo, err := NowGithubRepoCount(ctx, repoNm, token)
			if err != nil {
				log.Println(err)
				return
			}
			repos = append(repos, repo)
			log.Println(repoNm + " Start")
			stargazers := getStargazersCountByRepo(ctx, repoNm, token, repo)
			log.Println(repoNm + " DONE")
			lock.Lock()
			defer lock.Unlock()
			detaiRepos = append(detaiRepos, NewDetailsRepository(repo, stargazers))
		}(repoNm)
	}

	wg.Wait()
	gh := NewGitHub(repos, detaiRepos)
	return gh, nil
}

func NewDetailsRepository(repo GithubRepository, stargazers []Stargazer) ReadmeDetailsRepository {
	r := &ReadmeDetailsRepository{
		RepoName: repo.FullName,
		RepoURL:  repo.URL,
		StarCounts: map[string]int{
			"StarCount36MouthAgo": 0,
			"StarCount30MouthAgo": 0,
			"StarCount24MouthAgo": 0,
			"StarCount18MouthAgo": 0,
			"StarCount12MouthAgo": 0,
			"StarCount6MouthAgo":  0,
			"StarCountNow":        0,
		},
	}
	r.calculateStarCount(stargazers)
	r.RepoName = repo.FullName
	r.RepoURL = repo.URL
	return *r
}

func (r *ReadmeDetailsRepository) calculateStarCount(stargazers []Stargazer) {
	for _, star := range stargazers {
		r.updateStarCount("StarCountNow", star.StarredAt, 0)
		r.updateStarCount("StarCount6MouthAgo", star.StarredAt, -6)
		r.updateStarCount("StarCount12MouthAgo", star.StarredAt, -12)
		r.updateStarCount("StarCount18MouthAgo", star.StarredAt, -18)
		r.updateStarCount("StarCount24MouthAgo", star.StarredAt, -24)
		r.updateStarCount("StarCount30MouthAgo", star.StarredAt, -30)
		r.updateStarCount("StarCount36MouthAgo", star.StarredAt, -36)
	}
}

func (r *ReadmeDetailsRepository) updateStarCount(period string, starredAt time.Time, monthsAgo int) {
	var targetTime time.Time
	if monthsAgo == 0 {
		targetTime = time.Now().UTC()
	} else {
		targetTime = time.Now().UTC().AddDate(0, monthsAgo, 0)
	}

	if starredAt.Before(targetTime) {
		r.StarCounts[period]++
	}
}

func NowGithubRepoCount(ctx context.Context, name, token string) (GithubRepository, error) {
	url := baseURL + fmt.Sprintf("repos/%s", name)
	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return GithubRepository{}, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return GithubRepository{}, err
	}
	defer res.Body.Close()

	var repo GithubRepository
	if res.StatusCode == http.StatusOK {
		if err := json.Unmarshal(bts, &repo); err != nil {
			return GithubRepository{}, err
		}
	}

	return repo, nil
}

func getStargazersCountByRepo(ctx context.Context, name, token string, repo GithubRepository) []Stargazer {
	sem := make(chan bool, 4)
	var eg errgroup.Group
	var lock sync.Mutex
	var stargazers []Stargazer
	for page := 1; page <= lastPage(repo); page++ {
		sem <- true
		func(i int) {
			eg.Go(func() error {
				defer func() { <-sem }()
				result, err := getStargazersPage(ctx, repo, page, token)
				if errors.Is(err, ErrNoMorePages) {
					return err
				}
				if err != nil {
					return err
				}
				lock.Lock()
				defer lock.Unlock()
				stargazers = append(stargazers, result...)
				return nil
			})
		}(page)
		if err := eg.Wait(); err != nil {
			log.Println(err)
		}
	}

	return stargazers
}

func lastPage(repo GithubRepository) int {
	return totalPages(repo) + 1
}

func totalPages(repo GithubRepository) int {
	pageSize := 100
	return repo.StargazersCount / pageSize
}

func getStargazersPage(ctx context.Context, repo GithubRepository, page int, token string) ([]Stargazer, error) {
	var stars []Stargazer

	url := baseURL + fmt.Sprintf("repos/%s/stargazers?per_page=100&page=%d&", repo.FullName, page)
	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return stars, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(bts, &stars); err != nil {
			return nil, err
		}
		if len(stars) == 0 {
			return nil, fmt.Errorf("スターなし")
		}
		return stars, nil
	default:
		return nil, fmt.Errorf("その他のエラー")
	}
}

type GithubUser struct {
	AvatarURL string `json:"avatar_url"`
}

func GetRepoLogoUrl(repoName string, token string) (string, error) {
	owner := strings.Split(repoName, "/")[0]
	url := baseURL + fmt.Sprintf("users/%s", owner)
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

func GetRateLimit(token string) error {
	url := baseURL + rateLimit
	client := NewHttpClient(url, http.MethodGet, token)
	res, err := client.SendRequest()
	if err != nil {
		return err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var r map[string]interface{}
		if err := json.Unmarshal(bts, &r); err != nil {
			return err
		}
		fmt.Println(r)
		return nil
	case http.StatusNotModified:
		return fmt.Errorf("rate limit")
	case http.StatusNotFound:
		return fmt.Errorf("not Found")
	default:
		return fmt.Errorf("その他のエラー")
	}
}
