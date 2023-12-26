package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"star-golang-orms/domain/model"
	"time"
)

const (
	baseURL = "https://api.github.com/"
)

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
