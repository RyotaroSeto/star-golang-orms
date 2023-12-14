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

	if err = gh.MakeHTMLChartFile(); err != nil {
		return err
	}

	if err = model.ConvertHTMLToImage(); err != nil {
		return err
	}

	return gh.ReadmeEdit()
}

func (s *fetchService) createGitHub(ctx context.Context) (*model.GitHub, error) {
	if err := s.getGitHubRepos(ctx, s.gitHubRepo); err != nil {
		return nil, err
	}

	return &model.GitHub{
		Repositories:      s.repos,
		RepositoryDetails: s.detailRepos,
	}, nil
}

func (s *fetchService) getGitHubRepos(ctx context.Context, repo repository.GitHub) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-s.errCh:
		return err
	default:
		for _, repoNm := range model.TargetRepository {
			s.wg.Add(1)
			go func(repoNm string) {
				defer s.wg.Done()
				repo, err := s.gitHubRepo.GetRepository(ctx, model.RepositoryName(repoNm))
				if err != nil {
					s.errCh <- err
					return
				}
				s.repos = append(s.repos, *repo)

				log.Println(repoNm + " Start")
				stargazers := s.getStargazersCountByRepo(ctx, repo)
				log.Println(repoNm + " DONE")

				rd := model.NewRepositoryDetails(repo, stargazers)
				s.mu.Lock()
				defer s.mu.Unlock()
				s.detailRepos = append(s.detailRepos, rd)
			}(repoNm)
		}
	}
	s.wg.Wait()

	return nil
}

func (s *fetchService) getStargazersCountByRepo(ctx context.Context, repo *model.Repository) []model.Stargazer {
	var (
		stargazers = model.NewStargazers()
		wg         sync.WaitGroup
	)
	for page := 1; page <= model.LastPage(repo); page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			s.fetchStargazersPage(ctx, repo, page, stargazers)
		}(page)
	}
	wg.Wait()

	return stargazers.Stars
}

func (s *fetchService) fetchStargazersPage(ctx context.Context, repo *model.Repository, page int, stargazers *model.Stargazers) {
	pagers, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
	if errors.Is(err, domain.ErrNoMorePages) {
		return
	}
	if err != nil {
		s.errCh <- err
		return
	}

	stargazers.Add(*pagers)
}
