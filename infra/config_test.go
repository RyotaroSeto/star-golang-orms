package infra

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", githubToken)
	tests := []struct {
		name        string
		githubToken string
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name:        "success",
			githubToken: "github_token",
			assertion:   assert.NoError,
		},
		{
			name:        "failed",
			githubToken: "",
			assertion:   assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("GITHUB_TOKEN")

			if tt.githubToken != "" {
				os.Setenv("GITHUB_TOKEN", tt.githubToken)
			}
			tt.assertion(t, Load(context.Background()))
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		globalVar *Config
		want      *Config
	}{
		{
			name:      "success",
			globalVar: &Config{GitHubToken: "github_token"},
			want:      &Config{GitHubToken: "github_token"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = tt.globalVar
			assert.Equal(t, tt.want, Get())
		})
	}
}
