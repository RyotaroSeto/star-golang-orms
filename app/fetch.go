package app

import (
	"context"
	"log"
	"star-golang-orms/configs"
	"star-golang-orms/domain/repository"
	"star-golang-orms/domain/service"
	"star-golang-orms/pkg"
	"sync"
)

type fetchService struct {
	gitHubRepo repository.GitHub
	wg         sync.WaitGroup
	mu         sync.Mutex
	errCh      chan error
}

func NewFetchService(repo repository.GitHub) service.Fetcher {
	return &fetchService{
		gitHubRepo: repo,
		wg:         sync.WaitGroup{},
		mu:         sync.Mutex{},
		errCh:      make(chan error),
	}
}

func (s *fetchService) Start(ctx context.Context) error {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
		return err
	}

	gh, err := s.ExecGitHubAPI(ctx, config.GithubToken)
	if err != nil {
		log.Fatal("cannot exec github api", err)
		return err
	}

	err = gh.SortDesByStarCount()
	if err != nil {
		log.Fatal("cannot sort star count", err)
		return err
	}

	err = gh.MakeChart()
	if err != nil {
		log.Fatal("cannot make chart", err)
		return err
	}

	err = gh.Edit()
	if err != nil {
		log.Fatal("cannot edit readme", err)
		return err
	}
	return nil
}

func (s *fetchService) ExecGitHubAPI(ctx context.Context, token string) (pkg.GitHub, error) {
	var repos []pkg.GithubRepository
	var detaiRepos []pkg.ReadmeDetailsRepository

	wg := new(sync.WaitGroup)
	var lock sync.Mutex
	for _, repoNm := range pkg.TargetRepository {
		wg.Add(1)
		go func(repoNm string) {
			defer wg.Done()
			repo, err := pkg.NowGithubRepoCount(ctx, repoNm, token)
			// repo, err := s.gitHubRepo.GetRepository(ctx, model.RepositoryName(repoNm))
			if err != nil {
				log.Println(err)
				return
			}
			repos = append(repos, repo)
			log.Println(repoNm + " Start")
			// pagers, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
			stargazers := pkg.GetStargazersCountByRepo(ctx, token, repo)
			log.Println(repoNm + " DONE")
			lock.Lock()
			defer lock.Unlock()
			detaiRepos = append(detaiRepos, pkg.NewDetailsRepository(repo, stargazers))
		}(repoNm)
	}

	wg.Wait()
	gh := pkg.NewGitHub(repos, detaiRepos)
	return gh, nil
}
