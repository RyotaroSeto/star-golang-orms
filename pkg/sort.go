package pkg

import (
	"sort"
)

func (gh GitHub) SortDesByStarCount() error {
	err := githubRepositorySort(gh.GithubRepositorys)
	if err != nil {
		return err
	}
	err = readmeDetailsRepositorySort(gh.ReadmeDetailsRepositorys)
	if err != nil {
		return err
	}
	return nil
}

func githubRepositorySort(grs []GithubRepository) error {
	sort.Slice(grs, func(i, j int) bool {
		return grs[i].StargazersCount > grs[j].StargazersCount
	})
	return nil
}

func readmeDetailsRepositorySort(rds []ReadmeDetailsRepository) error {
	sort.Slice(rds, func(i, j int) bool {
		return rds[i].StarCountNow > rds[j].StarCountNow
	})
	return nil
}
