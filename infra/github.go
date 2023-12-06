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

var _ repository.GitHub = &GitHubRepository{}

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

func (r *GitHubRepository) newHttpRequest(ctx context.Context, rn model.RepositoryName) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+fmt.Sprintf("repos/%s", rn), nil)
	if err != nil {
		return nil, err
	}

	setRequestHeader(req)
	return req, nil
}

// httpmockを使ってテストを書く
func (r *GitHubRepository) GetRepository(ctx context.Context, rn model.RepositoryName) (*model.GitHubRepository, error) {
	req, err := r.newHttpRequest(ctx, rn)
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

	var repo *model.GitHubRepository
	if err := json.Unmarshal(b, &repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *GitHubRepository) GetStar(ctx context.Context) (int, error) {
	return 0, nil
}
