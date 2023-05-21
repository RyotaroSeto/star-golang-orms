package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"star-golang-orms/pkg"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	baseURL   = "https://api.github.com/"
	rateLimit = "rate_limit"
)

const (
	header = `# Golang ORMapper

| Project Name | Stars | Subscribers | Forks | Open Issues | Description | Create Update | Last Update |
| ------------ | ----- | ----------- | ----- | ----------- | ----------- | ----------- | ----------- |
`
	detailHeader = `
| 202111 | 202204 | 202301 | 202302 | 202303 |
| ------ | ------ | ------ | ------ | ------ |
`
)

type Stargazer struct {
	StarredAt string `json:"starred_at"`
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

type CheckMouth struct {
	RepoName        string
	RepoURL         string
	StarCount202111 int
	StarCount202204 int
	StarCount202301 int
	StarCount202302 int
	StarCount202303 int
}

func Edit(repos []GithubRepository, detaiRepos []CheckMouth) error {
	readme, err := os.Create("./README.md")
	if err != nil {
		return err
	}
	defer func() {
		_ = readme.Close()
	}()
	editREADME(readme, repos, detaiRepos)

	return nil
}

func editREADME(w io.Writer, repos []GithubRepository, detaiRepos []CheckMouth) {
	fmt.Fprint(w, header)
	for _, repo := range repos {
		fmt.Fprintf(w, "| [%s](%s) | %d | %d | %d | %d | %s | %v | %v |\n",
			repo.FullName,
			repo.URL,
			repo.StargazersCount,
			repo.SubscribersCount,
			repo.ForksCount,
			repo.OpenIssuesCount,
			repo.Description,
			repo.CreatedAt.Format("2006-01-02 15:04:05"),
			repo.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	for _, detaiRepo := range detaiRepos {
		fmt.Fprintf(w, "## [%s](%s)\n", detaiRepo.RepoName, detaiRepo.RepoURL)
		fmt.Fprint(w, detailHeader)
		fmt.Fprintf(w, "| %d | %d | %d | %d | %d |\n",
			detaiRepo.StarCount202111,
			detaiRepo.StarCount202204,
			detaiRepo.StarCount202301,
			detaiRepo.StarCount202302,
			detaiRepo.StarCount202303)
	}
}

func NowGithubRepoCount(ctx context.Context, name, token string) (GithubRepository, error) {
	url := baseURL + fmt.Sprintf("repos/%s", name)
	client := pkg.NewHttpClient(url, http.MethodGet, token)
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

func GetRepo(ctx context.Context, name, token string, repo GithubRepository) (CheckMouth, error) {
	sem := make(chan bool, 4)
	var eg errgroup.Group
	var lock sync.Mutex
	var stargazers []Stargazer
	for page := 1; page <= lastPage(repo); page++ {
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result, err := GetStargazersPage(ctx, repo, page, token)
			if errors.Is(err, pkg.ErrNoMorePages) {
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

	var cm CheckMouth
	cm.RepoName = strings.Split(name, "/")[1]
	cm.RepoURL = repo.URL
	for _, star := range stargazers {
		if star.StarredAt < "2021-11-01 00:00:00 +0000 UTC" {
			cm.StarCount202111++
		}
		if star.StarredAt < "2022-04-01 00:00:00 +0000 UTC" {
			cm.StarCount202204++
		}
		if star.StarredAt < "2023-01-01 00:00:00 +0000 UTC" {
			cm.StarCount202301++
		}
		if star.StarredAt < "2023-02-01 00:00:00 +0000 UTC" {
			cm.StarCount202302++
		}
		if star.StarredAt < "2023-03-01 00:00:00 +0000 UTC" {
			cm.StarCount202303++
		}
	}

	return cm, nil
}

// func GetRepo(ctx context.Context, name, token string) (GithubRepository, error) {
// 	repo, err := NowGithubRepoCount(name, token)
// 	if err != nil {
// 		log.Println(err)
// 		return GithubRepository{}, err
// 	}

// 	sem := make(chan bool, 4)
// 	var eg errgroup.Group
// 	var lock sync.Mutex
// 	var stargazers []Stargazer
// 	for page := 1; page <= lastPage(*repo); page++ {
// 		sem <- true
// 		page := page
// 		eg.Go(func() error {
// 			defer func() { <-sem }()
// 			result, err := GetStargazersPage(ctx, *repo, page, token)
// 			if errors.Is(err, pkg.ErrNoMorePages) {
// 				log.Println(err)
// 				return nil
// 			}
// 			if err != nil {
// 				log.Println(err)
// 				return err
// 			}
// 			lock.Lock()
// 			defer lock.Unlock()
// 			stargazers = append(stargazers, result...)
// 			return nil
// 		})
// 	}

// 	return *repo, nil

// 	// stargazers = append(stargazers, [{2017-07-07 02:50:15 +0000 UTC}])
// 	// log.Println(stargazers) //[{2017-07-07 02:50:15 +0000 UTC} {2017-07-07 05:06:33 +0000 UTC} {2017-07-07 10:56:49 +0000 UTC} {2017-07-07 11:25:36 +0000 UTC} {2017-07-07 19:42:38 +0000 UTC} {2017-07-08 01:06:01 +0000 UTC} ]

// 	// 時期によってのカウントの処理

// 	// READMEに書き込む処理
// }

func lastPage(repo GithubRepository) int {
	return totalPages(repo) + 1
}

func totalPages(repo GithubRepository) int {
	pageSize := 100
	return repo.StargazersCount / pageSize
}

func GetStargazersPage(ctx context.Context, repo GithubRepository, page int, token string) ([]Stargazer, error) {
	var stars []Stargazer

	url := baseURL + fmt.Sprintf("repos/%s/stargazers?per_page=100&page=%d&", repo.FullName, page)
	client := pkg.NewHttpClient(url, http.MethodGet, token)
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
	client := pkg.NewHttpClient(url, http.MethodGet, token)
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
	client := pkg.NewHttpClient(url, http.MethodGet, token)
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
