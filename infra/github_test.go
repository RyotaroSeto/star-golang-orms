package infra

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("GITHUB_TOKEN", "test_token")
	Load("../")
	os.Exit(m.Run())
}

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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type args struct {
		rn model.RepositoryName
	}
	tests := []struct {
		name       string
		args       args
		expReqJSON string
		respCode   int
		respBody   string
		want       *model.Repository
		assertion  assert.ErrorAssertionFunc
	}{
		// {
		// 	name: "success",
		// 	args: args{
		// 		rn: "test/test",
		// 	},
		// 	expReqJSON: ``,
		// 	respCode:   http.StatusOK,
		// 	respBody:   `{"full_name": "test/test", "html_url": "", "description": "test", "stargazers_count": 1, "subscribers_count": 1, "forks_count": 1, "open_issues_count": 1, "created_at": "2021-01-01T00:00:00Z", "updated_at": "2021-01-01T00:00:00Z"}`,
		// 	want: &model.Repository{
		// 		FullName:         "test/test",
		// 		URL:              "",
		// 		Description:      "test",
		// 		StargazersCount:  1,
		// 		SubscribersCount: 1,
		// 		ForksCount:       1,
		// 		OpenIssuesCount:  1,
		// 		CreatedAt:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		// 		UpdatedAt:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		// 	},
		// 	assertion: assert.NoError,
		// },
		{
			name: "failed. http status is BadRequest",
			args: args{
				rn: "test/test",
			},
			expReqJSON: ``,
			respCode:   http.StatusBadRequest,
			respBody:   `{"message": "test"}`,
			want:       nil,
			assertion:  assert.Error,
		},
		{
			name: "failed to unmarshal",
			args: args{
				rn: "test/test",
			},
			expReqJSON: ``,
			respCode:   http.StatusOK,
			respBody:   `{"full_name": "test/test", "html_url": "", "description": "test", "stargazers_count": "1", "subscribers_count": 1, "forks_count": 1, "open_issues_count": 1, "created_at": "2021-01-01T00:00:00Z", "updated_at": "2021-01-01T00:00:00Z"}`,
			want:       nil,
			assertion:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(http.MethodGet, baseURL+fmt.Sprintf("repos/%s", tt.args.rn),
				func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, r.Header.Get("Connection"), "keep-alive")
					assert.Equal(t, r.Header.Get("Authorization"), "token test_token")
					assert.Equal(t, r.Header.Get("Accept"), "application/vnd.github.v3.star+json")

					return httpmock.NewStringResponse(tt.respCode, tt.respBody), nil
				},
			)
			r := &GitHubRepository{
				client: &http.Client{},
			}
			got, err := r.GetRepository(context.Background(), tt.args.rn)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, 1, httpmock.GetCallCountInfo()[fmt.Sprintf("GET %s", baseURL+fmt.Sprintf("repos/%s", tt.args.rn))])
		})
	}
}

func TestGitHubRepository_GetStarPage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type args struct {
		repo *model.Repository
		page int
	}
	tests := []struct {
		name       string
		args       args
		expReqJSON string
		respCode   int
		respBody   string
		want       *[]model.Stargazer
		assertion  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				repo: &model.Repository{
					FullName: "test/test",
				},
				page: 1,
			},
			expReqJSON: ``,
			respCode:   http.StatusOK,
			respBody:   `[{"starred_at": "2021-01-01T00:00:00Z"}]`,
			want: &[]model.Stargazer{
				{StarredAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			assertion: assert.NoError,
		},
		{
			name: "failed. http status is BadRequest",
			args: args{
				repo: &model.Repository{
					FullName: "test/test",
				},
				page: 1,
			},
			expReqJSON: ``,
			respCode:   http.StatusBadRequest,
			respBody:   `{"message": "test"}`,
			want:       nil,
			assertion:  assert.Error,
		},
		{
			name: "failed to unmarshal",
			args: args{
				repo: &model.Repository{
					FullName: "test/test",
				},
				page: 1,
			},
			expReqJSON: ``,
			respCode:   http.StatusOK,
			respBody:   `{"starred_at": 1}`,
			want:       nil,
			assertion:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(http.MethodGet, baseURL+fmt.Sprintf("repos/%s/stargazers?per_page=100&page=%d&", tt.args.repo.FullName, tt.args.page),
				func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, r.Header.Get("Connection"), "keep-alive")
					assert.Equal(t, r.Header.Get("Authorization"), "token test_token")
					assert.Equal(t, r.Header.Get("Accept"), "application/vnd.github.v3.star+json")

					return httpmock.NewStringResponse(tt.respCode, tt.respBody), nil
				},
			)
			r := &GitHubRepository{
				client: &http.Client{},
			}

			got, err := r.GetStarPage(context.Background(), tt.args.repo, tt.args.page)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, 1, httpmock.GetCallCountInfo()[fmt.Sprintf("GET %s", baseURL+fmt.Sprintf("repos/%s/stargazers?per_page=100&page=%d&", tt.args.repo.FullName, tt.args.page))])
		})
	}
}
