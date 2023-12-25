package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"star-golang-orms/domain/model"
	"strings"
	"time"
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

func NewDetailsRepository(repo *model.Repository, stargazers []Stargazer) ReadmeDetailsRepository {
	r := &ReadmeDetailsRepository{
		RepoName: repo.FullName,
		RepoURL:  repo.URL,
		StarCounts: map[string]int{
			"StarCount72MouthAgo": 0,
			"StarCount60MouthAgo": 0,
			"StarCount48MouthAgo": 0,
			"StarCount36MouthAgo": 0,
			"StarCount24MouthAgo": 0,
			"StarCount12MouthAgo": 0,
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
		r.updateStarCount("StarCount12MouthAgo", star.StarredAt, -12)
		r.updateStarCount("StarCount24MouthAgo", star.StarredAt, -24)
		r.updateStarCount("StarCount36MouthAgo", star.StarredAt, -36)
		r.updateStarCount("StarCount48MouthAgo", star.StarredAt, -48)
		r.updateStarCount("StarCount60MouthAgo", star.StarredAt, -60)
		r.updateStarCount("StarCount72MouthAgo", star.StarredAt, -72)
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

func LastPage(repo *model.Repository) int {
	return totalPages(repo) + 1
}

func totalPages(repo *model.Repository) int {
	return repo.StargazersCount / 100
}

func GetStargazersPage(ctx context.Context, repo *model.Repository, page int, token string) ([]model.Stargazer, error) {
	var stars []model.Stargazer

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
			return nil, ErrNoStars
		}
		return stars, nil
	default:
		log.Println(res.StatusCode)
		return nil, ErrOtherReason
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
		return ErrRateLimit
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return ErrOtherReason
	}
}
