package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"star-golang-orms/pkg/errors"
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

func (r *GitHubRepository) newHttpRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	setRequestHeader(req)
	return req, nil
}

func (r *GitHubRepository) getFromGitHub(ctx context.Context, url string, result interface{}) (*http.Response, error) {
	req, err := r.newHttpRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req)
}

func (r *GitHubRepository) GetRepository(ctx context.Context, rn model.RepositoryName) (*model.Repository, error) {
	var repo model.Repository
	resp, err := r.getFromGitHub(ctx, baseURL+fmt.Sprintf("repos/%s", rn), &repo)
	if err != nil {
		return nil, errors.Newf(errors.InternalServerError, "failed to get repository: %s", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Newf(errors.InternalServerError, "failed to read response body: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(errors.InternalServerError, string(b))
	}

	if err := json.Unmarshal(b, &repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *GitHubRepository) GetStarPage(ctx context.Context, repo *model.Repository, page int) (*model.Stargazers, error) {
	var stars model.Stargazers
	resp, err := r.getFromGitHub(ctx, baseURL+fmt.Sprintf("repos/%s/stargazers?per_page=100&page=%d&", repo.FullName, page), &repo)
	if err != nil {
		return nil, errors.Newf(errors.InternalServerError, "failed to get repository: %s", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Newf(errors.InternalServerError, "failed to read response body: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.ErrOtherReason
	}

	if err := json.Unmarshal(b, &stars.Stars); err != nil {
		return nil, errors.Newf(errors.InternalServerError, "failed to unmarshal response body: %s", err)
	}

	if len(stars.Stars) == 0 {
		return nil, errors.ErrNoStars
	}

	return &stars, nil
}
