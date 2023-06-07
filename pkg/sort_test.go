package pkg

import (
	"testing"

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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
