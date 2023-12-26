package app

import (
	"context"
	"errors"
	"log"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"star-golang-orms/domain/service"
	"star-golang-orms/pkg"
	"sync"

	"golang.org/x/sync/errgroup"
)

type fetchService struct {
	gitHubRepo  repository.GitHub
	wg          sync.WaitGroup
	mu          sync.Mutex
	repos       model.Repositories
	detailRepos model.RepositoryDetails
	errCh       chan error
}

func NewFetchService(repo repository.GitHub) service.Fetcher {
	return &fetchService{
		gitHubRepo:  repo,
		wg:          sync.WaitGroup{},
		mu:          sync.Mutex{},
		repos:       make(model.Repositories, 0, len(model.TargetRepository)),
		detailRepos: make(model.RepositoryDetails, 0, len(model.TargetRepository)),
		errCh:       make(chan error),
	}
}

func (s *fetchService) Start(ctx context.Context) error {
	gh, err := s.createGitHub(ctx)
	if err != nil {
		return err
	}

	gh.ReadmeRepoAndDetailSort()
	if err = gh.MakeHTMLChartFile(); err != nil {
		return err
	}

	if err = model.ConvertHTMLToImage(); err != nil {
		return err
	}

	return gh.ReadmeEdit()
}

func (s *fetchService) createGitHub(ctx context.Context) (*model.GitHub, error) {
	wg := new(sync.WaitGroup)
	var lock sync.Mutex
	for _, repoNm := range model.TargetRepository {
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
			stargazers := s.getStargazersCountByRepo(ctx, repo)
			log.Println(repoNm + " DONE")

			rd := model.NewRepositoryDetails(repo, stargazers)
			lock.Lock()
			defer lock.Unlock()
			s.detailRepos = append(s.detailRepos, rd)
		}(repoNm)
	}

	wg.Wait()
	return &model.GitHub{
		Repositories:      s.repos,
		RepositoryDetails: s.detailRepos,
	}, nil
}

func (s *fetchService) getStargazersCountByRepo(ctx context.Context, repo *model.Repository) []model.Stargazer {
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
			if errors.Is(err, pkg.ErrNoMorePages) {
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

func (s *fetchService) fetchStargazersPage(ctx context.Context, repo *model.Repository, page int, stargazers *model.Stargazers) *[]model.Stargazer {
	pagers, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
	if errors.Is(err, pkg.ErrNoMorePages) {
		return nil
	}
	if err != nil {
		s.errCh <- err
		return nil
	}
	return pagers
}

// READMEの結果がソートされていない
