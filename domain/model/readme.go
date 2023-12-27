package model

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	header = `# Golang ORMapper Star 🎉🎉
The number of stars is expressed in an easy-to-understand manner for golang ormapper information with more than 1,000 stars. It can also display the number of stars at different times of the year.
If there are any other public repositories of golang orMapper, I'd be glad to hear about them!
`
	repoTable = `| No. | Project Name | Stars | Subscribers | Forks | Open Issues | Description | Createdate | Last Update |
| --- | ------------ | ----- | ----------- | ----- | ----------- | ----------- | ----------- | ----------- |
`
	divider = "|\n| --- | --- | --- | --- | --- | --- | --- |\n"
	README  = "README.md"
)

const (
	starCountZero = 0
	deployURL     = "https://ryotaroseto.github.io/star-golang-orms/output/orm_chart.html"
)

type GitHub struct {
	Repositories
	RepositoryDetails
}

type RepositoryDetail struct {
	RepoName   RepositoryName
	RepoURL    string
	StarCounts map[string]int
}

type RepositoryDetails []RepositoryDetail

type Repositories []Repository

func (gh GitHub) ReadmeEdit() error {
	readme, err := os.Create("./" + README)
	if err != nil {
		return err
	}
	defer func() {
		_ = readme.Close()
	}()
	gh.editREADME(readme)

	return nil
}

func NewRepositoryDetails(repo *Repository, stargazers []Stargazer) RepositoryDetail {
	r := &RepositoryDetail{
		RepoName: repo.RepositoryName(),
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

func (r RepositoryDetail) writeDetailRepo(w io.Writer) {
	repoHeader := fmt.Sprintf("## [%s](%s)\n", r.RepoName, r.RepoURL)
	fmt.Fprint(w, repoHeader)

	r.writeDetailRepoTable(w)
}

func (r RepositoryDetail) writeDetailRepoTable(w io.Writer) {
	fmt.Fprint(w, generateDetailRepoTableHeader())

	rowFormat := "| %s | %s | %s | %s | %s | %s | %s |\n"
	formattedStarCounts := r.formatStarCounts()
	fmt.Fprintf(w, rowFormat,
		formattedStarCounts["StarCount72MouthAgo"],
		formattedStarCounts["StarCount60MouthAgo"],
		formattedStarCounts["StarCount48MouthAgo"],
		formattedStarCounts["StarCount36MouthAgo"],
		formattedStarCounts["StarCount24MouthAgo"],
		formattedStarCounts["StarCount12MouthAgo"],
		formattedStarCounts["StarCountNow"])
}

func (r RepositoryDetail) formatStarCounts() map[string]string {
	formattedCounts := make(map[string]string)
	for period, count := range r.StarCounts {
		if count == starCountZero {
			formattedCounts[period] = "-"
		} else {
			formattedCounts[period] = strconv.Itoa(count)
		}
	}
	return formattedCounts
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

func (gh *GitHub) ReadmeRepoAndDetailSort() {
	gh.GithubRepositorySort()
	gh.ReadmeDetailsRepositorySort()
}

func (gh *GitHub) GithubRepositorySort() {
	tmpRepositories := make(Repositories, len(gh.Repositories))
	copy(tmpRepositories, gh.Repositories)
	sort.Sort(tmpRepositories)

	gh.Repositories = tmpRepositories
}

func (rs Repositories) Len() int {
	return len(rs)
}

func (rs Repositories) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Repositories) Less(i, j int) bool {
	return rs[i].StargazersCount > rs[j].StargazersCount
}

func (gh *GitHub) ReadmeDetailsRepositorySort() {
	tmpRepositoryDetails := make(RepositoryDetails, len(gh.RepositoryDetails))
	copy(tmpRepositoryDetails, gh.RepositoryDetails)
	sort.Sort(tmpRepositoryDetails)

	gh.RepositoryDetails = tmpRepositoryDetails
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

func LastPage(repo *Repository) int {
	return totalPages(repo) + 1
}

func totalPages(repo *Repository) int {
	return repo.StargazersCount / 100
}

func generateDetailRepoTableHeader() string {
	var detailHeader string

	dates := generateDates()
	for _, date := range dates {
		detailHeader += "| " + date + " "
	}
	detailHeader += divider

	return detailHeader
}

func generateDates() []string {
	now := time.Now()
	dates := make([]string, 7)

	for i := 0; i < len(dates); i++ {
		date := now.AddDate(0, -6*i, 0)
		dates[len(dates)-1-i] = date.Format(yyyymmddFormat)
	}

	return dates
}

func (gh GitHub) editREADME(w io.Writer) {
	writeHeader(w)
	writeChartJPEG(w)
	writeRepoTbl(w)
	writeRepositories(w, gh.Repositories)
	writeDetailRepositories(w, gh.RepositoryDetails)
}

func writeHeader(w io.Writer) {
	fmt.Fprint(w, header)
}

func writeChartJPEG(w io.Writer) {
	fmt.Fprintf(w, "[![Start数チャート](%s)](%s)\n", jpegFilePath, deployURL)
}

func writeRepoTbl(w io.Writer) {
	fmt.Fprint(w, repoTable)
}

func writeRepositories(w io.Writer, repos []Repository) {
	repoNo := 1
	for _, repo := range repos {
		repo.writeRepoRow(w, repoNo)
		repoNo++
	}
}

func writeDetailRepositories(w io.Writer, detailRepos []RepositoryDetail) {
	for _, d := range detailRepos {
		d.writeDetailRepo(w)
	}
}
