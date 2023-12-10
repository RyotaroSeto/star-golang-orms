package model

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
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

	README                     = "README.md"
	yyyymmddFormat             = time.DateOnly
	yyyymmddHHmmssHaihunFormat = time.DateTime
	starCountZero              = 0
	deployURL                  = "https://ryotaroseto.github.io/star-golang-orms/output/orm_chart.html"
)

const (
	htmlFilePath = "output/orm_chart.html"
	jpegFilePath = "output/orm_chart.jpeg"
)

type GitHub struct {
	Repositories
	RepositoryDetails
}

type RepositoryDetail struct {
	RepoName RepositoryName
	RepoURL  string
	// StarCounts IntervalStarCounts
	StarCounts map[string]int
}

type Repositories []Repository

func NewGitHub(repos Repositories, details RepositoryDetails) *GitHub {
	return &GitHub{
		Repositories:      repos,
		RepositoryDetails: details,
	}
}

func (rs *Repositories) AddRepo(repo *Repository) {
	*rs = append(*rs, *repo)
}

func (rd *RepositoryDetails) AddDetailRepo(repo *Repository, stargazers []Stargazer) {
	*rd = append(*rd, *NewRepositoryDetails(*repo, stargazers))
}

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

type RepositoryDetails []RepositoryDetail

func (rds RepositoryDetails) MakeHTMLChartFile() error {
	line := setUpChart()
	dates := generateDates()
	for _, d := range rds {
		line = d.makeLine(line, dates)
	}

	f, err := os.Create(htmlFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return line.Render(f)
}

func NewRepositoryDetails(repo Repository, stargazers []Stargazer) *RepositoryDetail {
	r := &RepositoryDetail{
		RepoName: repo.RepositoryName(),
		RepoURL:  repo.URL,
		// StarCounts: IntervalStarCounts{},
		// StarCounts: make(IntervalStarCounts, 7, len(stargazers)),

		StarCounts: map[string]int{
			"StarCount72MouthAgo": 0,
			"StarCount60MouthAgo": 0,
			"StarCount48MouthAgo": 0,
			"StarCount36MouthAgo": 0,
			"StarCount24MouthAgo": 0,
			"StarCountNow":        0,
			// "StarCountNow":        0,
			// "StarCount12MouthAgo": 0,
			// "StarCount24MouthAgo": 0,
			// "StarCount36MouthAgo": 0,
			// "StarCount48MouthAgo": 0,
			// "StarCount60MouthAgo": 0,
			// "StarCount72MouthAgo": 0,
		},
	}
	r.calculateStarCount(stargazers)
	return r
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

	// for _, v := range d.StarCounts {
	// 	starHistorys = append(starHistorys, opts.LineData{Value: v})
	// }

	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount72MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount60MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount48MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount36MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount24MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount12MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCountNow"]})

	return starHistorys
}

func GithubRepositorySort(rs Repositories) Repositories {
	tmpRepositories := make(Repositories, len(rs))
	copy(tmpRepositories, rs)
	sort.Sort(tmpRepositories)

	return tmpRepositories
}

func (rs Repositories) Len() int {
	return len(rs)
}

func (rs Repositories) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Repositories) Less(i, j int) bool {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].StargazersCount > rs[j].StargazersCount
	})
	return false
}

func ReadmeDetailsRepositorySort(rds RepositoryDetails) RepositoryDetails {
	tmpRepositoryDetails := make(RepositoryDetails, len(rds))
	copy(tmpRepositoryDetails, rds)
	sort.Sort(tmpRepositoryDetails)

	return tmpRepositoryDetails
}

func (rds RepositoryDetails) Len() int {
	return len(rds)
}

func (rds RepositoryDetails) Swap(i, j int) {
	rds[i], rds[j] = rds[j], rds[i]
}

func (rds RepositoryDetails) Less(i, j int) bool {
	sort.Slice(rds, func(i, j int) bool {
		return rds[i].StarCounts["StarCountNow"] > rds[j].StarCounts["StarCountNow"]
	})
	return false
}

func LastPage(repo Repository) int {
	return totalPages(repo) + 1
}

func totalPages(repo Repository) int {
	return repo.StargazersCount / 100
}

func ConvertHTMLToImage() error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)
	execAllocatorOpts := append(chromedp.DefaultExecAllocatorOptions[:], options...)
	ctx, cancel = chromedp.NewExecAllocator(ctx, execAllocatorOpts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	htmlContent, err := os.Open(htmlFilePath)
	if err != nil {
		return err
	}
	defer htmlContent.Close()

	data, err := io.ReadAll(htmlContent)
	if err != nil {
		return err
	}

	var buf []byte
	if err := chromedp.Run(ctx, screenshotTask(string(data), &buf)); err != nil {
		return err
	}

	return saveToFile(jpegFilePath, buf)
}

func screenshotTask(htmlContent string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html," + htmlContent),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Screenshot("body", res, chromedp.NodeVisible, chromedp.ByQuery),
	}
}

func saveToFile(filepath string, data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
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

func setUpChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
	)
	return line
}
