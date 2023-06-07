package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitHub(t *testing.T) {
	type args struct {
		gr []GithubRepository
		dr []ReadmeDetailsRepository
	}
	tests := []struct {
		name string
		args args
		want GitHub
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewGitHub(tt.args.gr, tt.args.dr))
		})
	}
}

func TestGitHub_Edit(t *testing.T) {
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
			tt.assertion(t, gh.Edit())
		})
	}
}
