package model

import (
	"fmt"
	"io"
	"sort"
	"time"
)

const (
	yyyymmddFormat             = time.DateOnly
	yyyymmddHHmmssHaihunFormat = time.DateTime
)

type Repository struct {
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

func (repo Repository) writeRowRepository(w io.Writer, repoNo int) {
	fmt.Fprintf(
		w,
		"| %d | [%s](%s) | %d | %d | %d | %d | %s | %s | %s |\n",
		repoNo,
		repo.FullName,
		repo.URL,
		repo.StargazersCount,
		repo.SubscribersCount,
		repo.ForksCount,
		repo.OpenIssuesCount,
		repo.Description,
		repo.CreatedAt.Format(yyyymmddHHmmssHaihunFormat),
		repo.UpdatedAt.Format(yyyymmddHHmmssHaihunFormat),
	)
}

func (repo *Repository) LastPage() int {
	return repo.totalPages() + 1
}

func (repo *Repository) totalPages() int {
	return repo.StargazersCount / 100
}

type RepositoryName string

func (r RepositoryName) String() string {
	return string(r)
}

type Repositories []Repository

func (rs *Repositories) GithubRepositorySort() {
	tmpRepositories := make(Repositories, len(*rs))
	copy(tmpRepositories, *rs)
	sort.Sort(tmpRepositories)

	*rs = tmpRepositories
}

func (rs Repositories) Len() int {
	return len(rs)
}

func (rs Repositories) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Repositories) Less(i, j int) bool {
	return rs[i].StargazersCount > rs[j].StargazersCount
}
