package model

import (
	"io"
	"os"
)

type GitHub struct {
	Repositories
	RepositoryDetails
}

func (gh *GitHub) ReadmeRepoAndDetailSort() {
	gh.GithubRepositorySort()
	gh.ReadmeDetailsRepositorySort()
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

func (gh GitHub) editREADME(w io.Writer) {
	writeHeader(w)
	writeChartJPEG(w)
	writeRepoTbl(w)
	writeRepositories(w, gh.Repositories)
	writeDetailRepositories(w, gh.RepositoryDetails)
}
