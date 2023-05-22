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

	divider = "|\n| --- | --- | --- | --- | --- | --- |\n"
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

type CheckMouth struct {
	RepoName            string
	RepoURL             string
	StarCount15MouthAgo int
	StarCount12MouthAgo int
	StarCount9MouthAgo  int
	StarCount6MouthAgo  int
	StarCount3MouthAgo  int
	StarCountNow        int
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

func editREADME(w io.Writer, repos []GithubRepository, detailRepos []CheckMouth) {
	writeHeader(w)
	writeRepositories(w, repos)
	writeDetailRepositories(w, detailRepos)
}

func writeHeader(w io.Writer) {
	fmt.Fprint(w, header)
}

func writeRepoRow(w io.Writer, repo GithubRepository) {
	rowFormat := "| [%s](%s) | %d | %d | %d | %d | %s | %s | %s |\n"
	createdAt := repo.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := repo.UpdatedAt.Format("2006-01-02 15:04:05")

	fmt.Fprintf(w, rowFormat, repo.FullName, repo.URL, repo.StargazersCount, repo.SubscribersCount, repo.ForksCount, repo.OpenIssuesCount, repo.Description, createdAt, updatedAt)
}

func writeRepositories(w io.Writer, repos []GithubRepository) {
	for _, repo := range repos {
		writeRepoRow(w, repo)
	}
}

func writeDetailRepositories(w io.Writer, detailRepos []CheckMouth) {

	for _, detailRepo := range detailRepos {
		writeDetailRepo(w, detailRepo)
	}
}

func writeDetailRepo(w io.Writer, detailRepo CheckMouth) {
	repoHeader := fmt.Sprintf("## [%s](%s)\n", detailRepo.RepoName, detailRepo.RepoURL)
	fmt.Fprint(w, repoHeader)

	writeDetailRepoTable(w, detailRepo)
}

func writeDetailRepoTable(w io.Writer, detailRepo CheckMouth) {
	fmt.Fprint(w, generateDetailRepoTableHeader())

	rowFormat := "| %d | %d | %d | %d | %d | %d |\n"
	fmt.Fprintf(w, rowFormat,
		detailRepo.StarCount15MouthAgo,
		detailRepo.StarCount12MouthAgo,
		detailRepo.StarCount9MouthAgo,
		detailRepo.StarCount6MouthAgo,
		detailRepo.StarCount3MouthAgo,
		detailRepo.StarCountNow)
}

func generateDetailRepoTableHeader() string {
	detailHeader := ""

	dates := generateDateHeaders()
	for _, date := range dates {
		detailHeader += "| " + date + " "
	}
	detailHeader += divider

	return detailHeader
}

func generateDateHeaders() []string {
	now := time.Now()
	layout := "20060102"
	dates := make([]string, 6)

	for i := 0; i < 6; i++ {
		date := now.AddDate(0, -3*i, 0)
		dates[i] = date.Format(layout)
	}

	return dates
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
		if star.StarredAt.Before(time.Now().UTC()) {
			cm.StarCountNow++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -3, 0)) {
			cm.StarCount3MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -6, 0)) {
			cm.StarCount6MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -9, 0)) {
			cm.StarCount9MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -12, 0)) {
			cm.StarCount12MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -15, 0)) {
			cm.StarCount15MouthAgo++
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
