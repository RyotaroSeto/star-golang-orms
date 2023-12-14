package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
)

const baseURL = "https://api.github.com/"

type GitHubRepository struct {
	client *http.Client
}

func NewGitHubRepository(ctx context.Context) repository.GitHub {
	return &GitHubRepository{
		client: registryHTTPClient(),
	}
}

func (r *GitHubRepository) get(ctx context.Context, req *http.Request) (*http.Response, error) {
	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *GitHubRepository) newHttpRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	setRequestHeader(req)
	return req, nil
}

func (r *GitHubRepository) GetRepository(ctx context.Context, rn model.RepositoryName) (*model.Repository, error) {
	req, err := r.newHttpRequest(ctx, baseURL+fmt.Sprintf("repos/%s", rn))
	if err != nil {
		return nil, err
	}

	resp, err := r.get(ctx, req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(string(b))
	}

	var repo *model.Repository
	if err := json.Unmarshal(b, &repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *GitHubRepository) GetStarPage(ctx context.Context, repo *model.Repository, page int) (*[]model.Stargazer, error) {
	req, err := r.newHttpRequest(ctx, baseURL+fmt.Sprintf("repos/%s/stargazers?per_page=100&page=%d&", repo.FullName, page))
	if err != nil {
		return nil, err
	}

	resp, err := r.get(ctx, req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(string(b))
	}

	var stars *[]model.Stargazer
	if err := json.Unmarshal(b, &stars); err != nil {
		return nil, err
	}

	return stars, nil
}
