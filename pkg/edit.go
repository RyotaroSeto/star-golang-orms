package pkg

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	header = `# Golang ORMapper Star
The number of stars is expressed in an easy-to-understand manner for golang ormapper information with more than 1,000 stars. It can also display the number of stars at different times of the year.
If there are any other public repositories of golang orMapper, I'd be glad to hear about them!

| Project Name | Stars | Subscribers | Forks | Open Issues | Description | Createdate | Last Update |
| ------------ | ----- | ----------- | ----- | ----------- | ----------- | ----------- | ----------- |
`

	divider = "|\n| --- | --- | --- | --- | --- | --- | --- |\n"

	README                     = "README.md"
	yyyymmddFormat             = "2006-01-02"
	yyyymmddHHmmssHaihunFormat = "2006-01-02 15:04:05"
	starCountZero              = 0
)

type GitHub struct {
	GithubRepositorys        []GithubRepository
	ReadmeDetailsRepositorys []ReadmeDetailsRepository
}

func NewGitHub(gr []GithubRepository, dr []ReadmeDetailsRepository) GitHub {
	return GitHub{
		GithubRepositorys:        gr,
		ReadmeDetailsRepositorys: dr,
	}
}

func (gh GitHub) Edit() error {
	readme, err := os.Create("./" + README)
	if err != nil {
		return err
	}
	defer func() {
		_ = readme.Close()
	}()
	editREADME(readme, gh.GithubRepositorys, gh.ReadmeDetailsRepositorys)

	return nil
}

func editREADME(w io.Writer, repos []GithubRepository, detailRepos []ReadmeDetailsRepository) {
	writeHeader(w)
	writeRepositories(w, repos)
	writeDetailRepositories(w, detailRepos)
}

func writeHeader(w io.Writer) {
	fmt.Fprint(w, header)
}

func writeRepositories(w io.Writer, repos []GithubRepository) {
	for _, repo := range repos {
		repo.writeRepoRow(w)
	}
}

func (repo GithubRepository) writeRepoRow(w io.Writer) {
	rowFormat := "| [%s](%s) | %d | %d | %d | %d | %s | %s | %s |\n"
	createdAt := repo.CreatedAt.Format(yyyymmddHHmmssHaihunFormat)
	updatedAt := repo.UpdatedAt.Format(yyyymmddHHmmssHaihunFormat)

	fmt.Fprintf(w, rowFormat, repo.FullName, repo.URL, repo.StargazersCount, repo.SubscribersCount, repo.ForksCount, repo.OpenIssuesCount, repo.Description, createdAt, updatedAt)
}

func writeDetailRepositories(w io.Writer, detailRepos []ReadmeDetailsRepository) {
	for _, d := range detailRepos {
		d.writeDetailRepo(w)
	}
}

func (r ReadmeDetailsRepository) writeDetailRepo(w io.Writer) {
	repoHeader := fmt.Sprintf("## [%s](%s)\n", r.RepoName, r.RepoURL)
	fmt.Fprint(w, repoHeader)

	r.writeDetailRepoTable(w)
}

func (r ReadmeDetailsRepository) writeDetailRepoTable(w io.Writer) {
	fmt.Fprint(w, generateDetailRepoTableHeader())

	rowFormat := "| %s | %s | %s | %s | %s | %s | %s |\n"
	formattedStarCounts := r.formatStarCounts()
	fmt.Fprintf(w, rowFormat,
		formattedStarCounts["StarCount36MouthAgo"],
		formattedStarCounts["StarCount30MouthAgo"],
		formattedStarCounts["StarCount24MouthAgo"],
		formattedStarCounts["StarCount18MouthAgo"],
		formattedStarCounts["StarCount12MouthAgo"],
		formattedStarCounts["StarCount6MouthAgo"],
		formattedStarCounts["StarCountNow"])
}

func (r ReadmeDetailsRepository) formatStarCounts() map[string]string {
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

func generateDetailRepoTableHeader() string {
	detailHeader := ""

	dates := generateDateHeaders()
	for _, date := range dates {
		detailHeader += "| " + date + " "
	}
	detailHeader += divider

	return detailHeader
}

func generateDateHeaders() []string {
	now := time.Now()
	dates := make([]string, 7)

	for i := 0; i < len(dates); i++ {
		date := now.AddDate(0, -6*i, 0)
		dates[len(dates)-1-i] = date.Format(yyyymmddFormat)
	}

	return dates
}
