package app

import (
	"context"
	"log"
	"star-golang-orms/configs"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"star-golang-orms/domain/service"
	"star-golang-orms/pkg"
	"sync"
)

type fetchService struct {
	gitHubRepo repository.GitHub
	wg         sync.WaitGroup
	mu         sync.Mutex
	repos      model.Repositories
	errCh      chan error
}

func NewFetchService(repo repository.GitHub) service.Fetcher {
	return &fetchService{
		gitHubRepo: repo,
		wg:         sync.WaitGroup{},
		mu:         sync.Mutex{},
		repos:      make(model.Repositories, 0, len(model.TargetRepository)),
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
	var detaiRepos []pkg.ReadmeDetailsRepository

	wg := new(sync.WaitGroup)
	var lock sync.Mutex
	for _, repoNm := range pkg.TargetRepository {
		wg.Add(1)
		go func(repoNm string) {
			defer wg.Done()
			repo, err := s.gitHubRepo.GetRepository(ctx, model.RepositoryName(repoNm))
			if err != nil {
				log.Println(err)
				return
			}
			s.repos = append(s.repos, *repo)
			log.Println(repoNm + " Start")
			stargazers := pkg.GetStargazersCountByRepo(ctx, token, repo)
			// stargazers := s.GetStargazersCountByRepo(ctx, token, repo)
			log.Println(repoNm + " DONE")
			lock.Lock()
			defer lock.Unlock()
			detaiRepos = append(detaiRepos, pkg.NewDetailsRepository(repo, stargazers))
		}(repoNm)
	}

	wg.Wait()
	gh := pkg.NewGitHub(s.repos, detaiRepos)
	return gh, nil
}

// func (s *fetchService) GetStargazersCountByRepo(ctx context.Context, token string, repo pkg.GithubRepository) []pkg.Stargazer {
// 	sem := make(chan bool, 4)
// 	var eg errgroup.Group
// 	var lock sync.Mutex
// 	var stargazers []pkg.Stargazer
// 	for page := 1; page <= pkg.LastPage(repo); page++ {
// 		sem <- true
// 		page := page
// 		eg.Go(func() error {
// 			defer func() { <-sem }()
// 			pagers, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
// 			// if errors.Is(err, ErrNoMorePages) {
// 			// 	return nil
// 			// }
// 			if err != nil {
// 				return err
// 			}
// 			lock.Lock()
// 			defer lock.Unlock()
// 			stargazers = append(stargazers, pagers...)
// 			return nil
// 		})
// 	}
// 	if err := eg.Wait(); err != nil {
// 		log.Println(err)
// 	}

// 	return stargazers
// }
