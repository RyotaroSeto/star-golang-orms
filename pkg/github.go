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
	RepoName            string
	RepoURL             string
	StarCount30MouthAgo int
	StarCount24MouthAgo int
	StarCount18MouthAgo int
	StarCount12MouthAgo int
	StarCount6MouthAgo  int
	StarCountNow        int
}

func NewDetailsRepository(repoName, repoURL string, stargazers []Stargazer) *ReadmeDetailsRepository {
	var r ReadmeDetailsRepository
	r.calculateStarCount(stargazers)
	r.RepoName = strings.Split(repoName, "/")[1]
	r.RepoURL = repoURL
	return &r
}

func (r *ReadmeDetailsRepository) calculateStarCount(stargazers []Stargazer) {
	for _, star := range stargazers {
		if star.StarredAt.Before(time.Now().UTC()) {
			r.StarCountNow++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -6, 0)) {
			r.StarCount6MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -12, 0)) {
			r.StarCount12MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -18, 0)) {
			r.StarCount18MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -24, 0)) {
			r.StarCount24MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -30, 0)) {
			r.StarCount30MouthAgo++
		}
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

func GetRepo(ctx context.Context, name, token string, repo GithubRepository) (ReadmeDetailsRepository, error) {
	sem := make(chan bool, 4)
	var eg errgroup.Group
	var lock sync.Mutex
	var stargazers []Stargazer
	for page := 1; page <= lastPage(repo); page++ {
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result, err := getStargazersPage(ctx, repo, page, token)
			if errors.Is(err, ErrNoMorePages) {
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

	detailsRepository := NewDetailsRepository(name, repo.URL, stargazers)
	return *detailsRepository, nil
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
