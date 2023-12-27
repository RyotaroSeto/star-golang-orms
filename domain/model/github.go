package model

import (
	"io"
	"os"
)

type GitHub struct {
	Repositories
	RepositoryDetails
}

func (gh *GitHub) RepoAndDetailAscSort() {
	gh.githubRepositorySort()
	gh.readmeDetailsRepositorySort()
}

func (gh GitHub) EditREADME() (err error) {
	readme, err := os.Create("./" + README)
	if err != nil {
		return
	}
	defer func() {
		_ = readme.Close()
	}()
	gh.editStart(readme)

	return
}

func (gh GitHub) editStart(w io.Writer) {
	writeHeader(w)
	writeChartJPEG(w)
	writeRepoTbl(w)
	writeRepositories(w, gh.Repositories)
	writeDetailRepositories(w, gh.RepositoryDetails)
}
