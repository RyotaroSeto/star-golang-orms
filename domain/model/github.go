package model

import (
	"fmt"
	"io"
	"sync"
	"time"
)

const (
	yyyymmddFormat             = time.DateOnly
	yyyymmddHHmmssHaihunFormat = time.DateTime
)

type Stargazer struct {
	StarredAt time.Time `json:"starred_at"`
}

type Stargazers struct {
	Stars []Stargazer
	lock  sync.Mutex
}

func NewStargazers() *Stargazers {
	return &Stargazers{
		Stars: make([]Stargazer, 0),
		lock:  sync.Mutex{},
	}
}

func (ss *Stargazers) Add(stargazers []Stargazer) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.Stars = append(ss.Stars, stargazers...)
}

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

func (r Repository) RepositoryName() RepositoryName {
	return RepositoryName(r.FullName)
}

func (repo Repository) writeRepoRow(w io.Writer, repoNo int) {
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

type RepositoryName string

func (r RepositoryName) String() string {
	return string(r)
}
