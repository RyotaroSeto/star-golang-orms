package model

import (
	"fmt"
	"io"
	"time"
)

type Stargazer struct {
	StarredAt time.Time `json:"starred_at"`
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
	rowFormat := "| %d | [%s](%s) | %d | %d | %d | %d | %s | %s | %s |\n"
	createdAt := repo.CreatedAt.Format(yyyymmddHHmmssHaihunFormat)
	updatedAt := repo.UpdatedAt.Format(yyyymmddHHmmssHaihunFormat)

	fmt.Fprintf(w, rowFormat, repoNo, repo.FullName, repo.URL, repo.StargazersCount, repo.SubscribersCount, repo.ForksCount, repo.OpenIssuesCount, repo.Description, createdAt, updatedAt)
}

type RepositoryName string

func (r RepositoryName) String() string {
	return string(r)
}

type IntervalStarCount struct {
	StarCount72MouthAgo int
	StarCount60MouthAgo int
	StarCount48MouthAgo int
	StarCount36MouthAgo int
	StarCount24MouthAgo int
	StarCount12MouthAgo int
	StarCountNow        int
}

type IntervalStarCounts []IntervalStarCount
