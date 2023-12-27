package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepository_LastPage(t *testing.T) {
	type fields struct {
		FullName         string
		URL              string
		Description      string
		StargazersCount  int
		SubscribersCount int
		ForksCount       int
		OpenIssuesCount  int
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				StargazersCount: 100,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				FullName:         tt.fields.FullName,
				URL:              tt.fields.URL,
				Description:      tt.fields.Description,
				StargazersCount:  tt.fields.StargazersCount,
				SubscribersCount: tt.fields.SubscribersCount,
				ForksCount:       tt.fields.ForksCount,
				OpenIssuesCount:  tt.fields.OpenIssuesCount,
				CreatedAt:        tt.fields.CreatedAt,
				UpdatedAt:        tt.fields.UpdatedAt,
			}
			assert.Equal(t, tt.want, repo.LastPage())
		})
	}
}

func TestRepositoryName_String(t *testing.T) {
	tests := []struct {
		name string
		r    RepositoryName
		want string
	}{
		{
			name: "success",
			r:    "test",
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.r.String())
		})
	}
}
