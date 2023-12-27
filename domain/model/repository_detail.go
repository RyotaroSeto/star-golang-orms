package model

import (
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type RepositoryDetail struct {
	RepoName   RepositoryName
	RepoURL    string
	StarCounts map[string]int
}

func NewRepositoryDetails(repo *Repository, stargazers []Stargazer) RepositoryDetail {
	r := &RepositoryDetail{
		RepoName: RepositoryName(repo.FullName),
		RepoURL:  repo.URL,
		StarCounts: map[string]int{
			"StarCount72MouthAgo": 0,
			"StarCount60MouthAgo": 0,
			"StarCount48MouthAgo": 0,
			"StarCount36MouthAgo": 0,
			"StarCount24MouthAgo": 0,
			"StarCountNow":        0,
		},
	}
	r.calculateStarCount(stargazers)

	return *r
}

func (r *RepositoryDetail) calculateStarCount(stargazers []Stargazer) {
	for _, star := range stargazers {
		r.updateStarCount("StarCountNow", star.StarredAt, 0)
		r.updateStarCount("StarCount12MouthAgo", star.StarredAt, -12)
		r.updateStarCount("StarCount24MouthAgo", star.StarredAt, -24)
		r.updateStarCount("StarCount36MouthAgo", star.StarredAt, -36)
		r.updateStarCount("StarCount48MouthAgo", star.StarredAt, -48)
		r.updateStarCount("StarCount60MouthAgo", star.StarredAt, -60)
		r.updateStarCount("StarCount72MouthAgo", star.StarredAt, -72)
	}
}

func (r *RepositoryDetail) updateStarCount(period string, starredAt time.Time, monthsAgo int) {
	var targetTime time.Time
	if monthsAgo == 0 {
		targetTime = time.Now().UTC()
	} else {
		targetTime = time.Now().UTC().AddDate(0, monthsAgo, 0)
	}

	if starredAt.Before(targetTime) {
		r.StarCounts[period]++
	}
}

func (d RepositoryDetail) makeLine(line *charts.Line, dates []string) *charts.Line {
	starHistorys := d.generateStarHistorys()
	line.SetXAxis(dates).AddSeries(d.RepoName.String(), starHistorys)
	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line
}

func (d RepositoryDetail) generateStarHistorys() []opts.LineData {
	starHistorys := make([]opts.LineData, 0, 7)
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount72MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount60MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount48MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount36MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount24MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount12MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCountNow"]})

	return starHistorys
}

type RepositoryDetails []RepositoryDetail

func (rds *RepositoryDetails) ReadmeDetailsRepositorySort() {
	tmpRepositoryDetails := make(RepositoryDetails, len(*rds))
	copy(tmpRepositoryDetails, *rds)
	sort.Sort(tmpRepositoryDetails)

	*rds = tmpRepositoryDetails
}

func (rds RepositoryDetails) Len() int {
	return len(rds)
}

func (rds RepositoryDetails) Swap(i, j int) {
	rds[i], rds[j] = rds[j], rds[i]
}

func (rds RepositoryDetails) Less(i, j int) bool {
	return rds[i].StarCounts["StarCountNow"] > rds[j].StarCounts["StarCountNow"]
}
