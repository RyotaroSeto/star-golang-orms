package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
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
		repo.writeRowRepository(w, repoNo)
		repoNo++
	}
}

func writeDetailRepositories(w io.Writer, detailRepos []RepositoryDetail) {
	for _, d := range detailRepos {
		writeDetailRepo(w, d)
	}
}

func writeDetailRepo(w io.Writer, rd RepositoryDetail) {
	repoHeader := fmt.Sprintf("## [%s](%s)\n", rd.RepoName, rd.RepoURL)
	fmt.Fprint(w, repoHeader)

	writeDetailRepoTable(w, rd)
}

func writeDetailRepoTable(w io.Writer, rd RepositoryDetail) {
	fmt.Fprint(w, generateDetailRepoTableHeader())

	formattedStarCounts := formatStarCounts(rd)
	fmt.Fprintf(
		w,
		"| %s | %s | %s | %s | %s | %s | %s |\n",
		formattedStarCounts["StarCount72MouthAgo"],
		formattedStarCounts["StarCount60MouthAgo"],
		formattedStarCounts["StarCount48MouthAgo"],
		formattedStarCounts["StarCount36MouthAgo"],
		formattedStarCounts["StarCount24MouthAgo"],
		formattedStarCounts["StarCount12MouthAgo"],
		formattedStarCounts["StarCountNow"],
	)
}

func formatStarCounts(rd RepositoryDetail) map[string]string {
	formattedCounts := make(map[string]string)
	for period, count := range rd.StarCounts {
		if count == starCountZero {
			formattedCounts[period] = "-"
		} else {
			formattedCounts[period] = strconv.Itoa(count)
		}
	}
	return formattedCounts
}
