package app

import (
	"context"
	"errors"
	"log"
	"star-golang-orms/domain"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"star-golang-orms/domain/service"
	"sync"

	"golang.org/x/sync/errgroup"
)

type fetchService struct {
	gitHubRepo repository.GitHub
	mu         sync.Mutex
}

func NewFetchService(repo repository.GitHub) service.Fetcher {
	return &fetchService{
		gitHubRepo: repo,
		mu:         sync.Mutex{},
	}
}

func (s *fetchService) Start(ctx context.Context) error {
	gh, err := s.createGitHub(ctx)
	if err != nil {
		return err
	}

	err = gh.MakeHTMLChartFile()
	if err != nil {
		return err
	}

	err = model.ConvertHTMLToImage()
	if err != nil {
		return err
	}

	err = gh.ReadmeEdit()
	if err != nil {
		return err
	}

	return nil
}

func (s *fetchService) createGitHub(ctx context.Context) (*model.GitHub, error) {
	var repos model.Repositories
	var detaiRepos model.RepositoryDetails

	wg := new(sync.WaitGroup)
	for _, repoNm := range model.TargetRepository {
		wg.Add(1)
		go func(repoNm string) {
			defer wg.Done()
			repo, err := s.gitHubRepo.GetRepository(ctx, model.RepositoryName(repoNm))
			if err != nil {
				log.Println(err)
				return
			}
			repos = append(repos, *repo)

			log.Println(repoNm + " Start")
			stargazers := s.getStargazersCountByRepo(ctx, *repo)
			log.Println(repoNm + " DONE")

			s.mu.Lock()
			defer s.mu.Unlock()
			detaiRepos = append(detaiRepos, *model.NewRepositoryDetails(*repo, stargazers))
		}(repoNm)
	}
	wg.Wait()

	return model.NewGitHub(
		model.GithubRepositorySort(repos),
		model.ReadmeDetailsRepositorySort(detaiRepos),
	), nil
}

// defer wg.Done()
// 			repo, err := s.gitHubRepo.GetRepository(ctx, model.RepositoryName(repoNm))
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			rs.AddRepo(repo)

// 			log.Println(repoNm + " Start")
// 			stargazers := s.getStargazersCountByRepo(ctx, *repo)
// 			log.Println(repoNm + " DONE")

// s.mu.Lock()
// defer s.mu.Unlock()
// dr.AddDetailRepo(repo, stargazers)
func (s *fetchService) getStargazersCountByRepo(ctx context.Context, repo model.Repository) []model.Stargazer {
	sem := make(chan bool, 4)
	var eg errgroup.Group
	var lock sync.Mutex
	var stargazers []model.Stargazer
	for page := 1; page <= model.LastPage(repo); page++ {
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
			if errors.Is(err, domain.ErrNoMorePages) {
				return nil
			}
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			stargazers = append(stargazers, *result...)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		log.Println(err)
	}

	return stargazers
}
