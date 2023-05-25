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
	header    = `# Golang ORMapper Star
Express information on golang ormapper in a clear manner. It also displays the number of stars at different times of the year.

| Project Name | Stars | Subscribers | Forks | Open Issues | Description | Createdate | Last Update |
| ------------ | ----- | ----------- | ----- | ----------- | ----------- | ----------- | ----------- |
`

	divider = "|\n| --- | --- | --- | --- | --- | --- |\n"

	README                     = "README.md"
	yyyymmddFormat             = "20060102"
	yyyymmddHHmmssHaihunFormat = "2006-01-02 15:04:05"
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
	StarCount15MouthAgo int
	StarCount12MouthAgo int
	StarCount9MouthAgo  int
	StarCount6MouthAgo  int
	StarCount3MouthAgo  int
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
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -3, 0)) {
			r.StarCount3MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -6, 0)) {
			r.StarCount6MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -9, 0)) {
			r.StarCount9MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -12, 0)) {
			r.StarCount12MouthAgo++
		}
		if star.StarredAt.Before(time.Now().UTC().AddDate(0, -15, 0)) {
			r.StarCount15MouthAgo++
		}
	}
}

func Edit(repos []GithubRepository, detaiRepos []ReadmeDetailsRepository) error {
	readme, err := os.Create("./" + README)
	if err != nil {
		return err
	}
	defer func() {
		_ = readme.Close()
	}()
	editREADME(readme, repos, detaiRepos)

	return nil
}

func editREADME(w io.Writer, repos []GithubRepository, detailRepos []ReadmeDetailsRepository) {
	writeHeader(w)
	writeRepositories(w, repos)
	writeDetailRepositories(w, detailRepos)
}

func writeHeader(w io.Writer) {
	fmt.Fprint(w, header)
}

func writeRepoRow(w io.Writer, repo GithubRepository) {
	rowFormat := "| [%s](%s) | %d | %d | %d | %d | %s | %s | %s |\n"
	createdAt := repo.CreatedAt.Format(yyyymmddHHmmssHaihunFormat)
	updatedAt := repo.UpdatedAt.Format(yyyymmddHHmmssHaihunFormat)

	fmt.Fprintf(w, rowFormat, repo.FullName, repo.URL, repo.StargazersCount, repo.SubscribersCount, repo.ForksCount, repo.OpenIssuesCount, repo.Description, createdAt, updatedAt)
}

func writeRepositories(w io.Writer, repos []GithubRepository) {
	for _, repo := range repos {
		writeRepoRow(w, repo)
	}
}

func writeDetailRepositories(w io.Writer, detailRepos []ReadmeDetailsRepository) {
	for _, d := range detailRepos {
		d.writeDetailRepo(w)
	}
}

func (r ReadmeDetailsRepository) writeDetailRepo(w io.Writer) {
	repoHeader := fmt.Sprintf("## [%s](%s)\n", r.RepoName, r.RepoURL)
	fmt.Fprint(w, repoHeader)

	r.writeDetailRepoTable(w)
}

func (r ReadmeDetailsRepository) writeDetailRepoTable(w io.Writer) {
	fmt.Fprint(w, generateDetailRepoTableHeader())

	rowFormat := "| %d | %d | %d | %d | %d | %d |\n"
	fmt.Fprintf(w, rowFormat,
		r.StarCount15MouthAgo,
		r.StarCount12MouthAgo,
		r.StarCount9MouthAgo,
		r.StarCount6MouthAgo,
		r.StarCount3MouthAgo,
		r.StarCountNow)
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
	dates := make([]string, 6)

	for i := 0; i < len(dates); i++ {
		date := now.AddDate(0, -3*i, 0)
		dates[len(dates)-1-i] = date.Format("20060102")
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
