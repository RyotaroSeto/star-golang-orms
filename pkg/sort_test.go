package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitHub_SortDesByStarCount(t *testing.T) {
	type fields struct {
		GithubRepositorys        []GithubRepository
		ReadmeDetailsRepositorys []ReadmeDetailsRepository
	}
	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "正常系",
			fields: fields{
				GithubRepositorys: []GithubRepository{
					{
						FullName:         "test",
						URL:              "test",
						Description:      "test",
						StargazersCount:  1,
						SubscribersCount: 1,
						ForksCount:       1,
						OpenIssuesCount:  1,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
				},
				ReadmeDetailsRepositorys: []ReadmeDetailsRepository{
					{
						RepoName:   "test",
						RepoURL:    "test",
						StarCounts: map[string]int{"2020-01-01": 1},
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "異常系",
			fields: fields{
				GithubRepositorys: []GithubRepository{
					{
						FullName:         "test",
						URL:              "test",
						Description:      "test",
						StargazersCount:  1,
						SubscribersCount: 1,
						ForksCount:       1,
						OpenIssuesCount:  1,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
				},
				ReadmeDetailsRepositorys: []ReadmeDetailsRepository{
					{
						RepoName:   "test",
						RepoURL:    "test",
						StarCounts: map[string]int{"2020-01-01": 1},
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gh := GitHub{
				GithubRepositorys:        tt.fields.GithubRepositorys,
				ReadmeDetailsRepositorys: tt.fields.ReadmeDetailsRepositorys,
			}
			tt.assertion(t, gh.SortDesByStarCount())
		})
	}
}

func Test_githubRepositorySort(t *testing.T) {
	type args struct {
		grs []GithubRepository
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "正常系",
			args: args{
				grs: []GithubRepository{
					{
						FullName:         "test",
						URL:              "test",
						Description:      "test",
						StargazersCount:  1,
						SubscribersCount: 1,
						ForksCount:       1,
						OpenIssuesCount:  1,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "異常系",
			args: args{
				grs: []GithubRepository{
					{
						FullName:         "test",
						URL:              "test",
						Description:      "test",
						StargazersCount:  1,
						SubscribersCount: 1,
						ForksCount:       1,
						OpenIssuesCount:  1,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, githubRepositorySort(tt.args.grs))
		})
	}
}

func Test_readmeDetailsRepositorySort(t *testing.T) {
	type args struct {
		rds []ReadmeDetailsRepository
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, readmeDetailsRepositorySort(tt.args.rds))
		})
	}
}
