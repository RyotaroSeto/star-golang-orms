package infra

import (
	"context"
	"io"
	"net/http"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewGitHubRepository(t *testing.T) {
	tests := []struct {
		name string
		want repository.GitHub
	}{
		{
			name: "success",
			want: &GitHubRepository{
				client: &http.Client{
					Timeout: time.Duration(1000) * time.Second,
					Transport: &http.Transport{
						IdleConnTimeout: time.Duration(180) * time.Second,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewGitHubRepository(context.Background()))
		})
	}
}

func TestGitHubRepository_GetRepository(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "test_token")
	Load(".")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type args struct {
		rn model.RepositoryName
	}
	tests := []struct {
		name      string
		args      args
		want      *model.GitHubRepository
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(http.MethodGet,
				func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, r.Header.Get("Connection"), "keep-alive")
					assert.Equal(t, r.Header.Get("Authorization"), "token test_token")
					assert.Equal(t, r.Header.Get("Accept"), "application/vnd.github.v3.star+json")

					b, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatal(err)
					}
					defer r.Body.Close()
					assert.JSONEq(t, tt.args.rn, string(b))

					return httpmock.NewStringResponse(tt.respCode, tt.respBody), nil
				})

			r := &GitHubRepository{
				client: &http.Client{},
			}
			got, err := r.GetRepository(context.Background(), tt.args.rn)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGitHubRepository_GetStarPage(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx  context.Context
		repo model.GitHubRepository
		page int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *model.Stargazer
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GitHubRepository{
				client: tt.fields.client,
			}
			got, err := r.GetStarPage(tt.args.ctx, tt.args.repo, tt.args.page)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
