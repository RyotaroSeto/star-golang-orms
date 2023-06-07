package pkg

import (
	"testing"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/stretchr/testify/assert"
)

func TestGitHub_MakeChart(t *testing.T) {
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
			tt.assertion(t, gh.MakeChart())
		})
	}
}

func TestGitHub_makeHTMLChartFile(t *testing.T) {
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
			tt.assertion(t, gh.makeHTMLChartFile())
		})
	}
}

func Test_setUpChart(t *testing.T) {
	tests := []struct {
		name string
		want *charts.Line
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, setUpChart())
		})
	}
}
