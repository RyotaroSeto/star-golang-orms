package app

import (
	"context"
	"errors"
	"log"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"star-golang-orms/domain/service"
	starError "star-golang-orms/pkg/errors"
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

	gh.RepoAndDetailAscSort()
	if err = gh.MakeHTMLChartFile(); err != nil {
		return starError.Newf(starError.InternalServerError, "failed to create html file: %s", err)
	}

	if err = model.ConvertHTMLToImage(); err != nil {
		return starError.Newf(starError.InternalServerError, "failed to convert html to image: %s", err)
	}

	return gh.EditREADME()
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
		sem        = make(chan bool, 4)
		eg         errgroup.Group
		stargazers = model.NewStargazers()
	)
	for page := 1; page <= repo.LastPage(); page++ {
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result := s.fetchStargazersPage(ctx, repo, page, stargazers)
			stargazers.Add(result.Stars)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		s.errCh <- err
	}

	return stargazers.Stars
}

func (s *fetchService) fetchStargazersPage(ctx context.Context, repo *model.Repository, page int, stargazers *model.Stargazers) *model.Stargazers {
	pagers, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
	if errors.Is(err, starError.ErrNoMorePages) {
		return nil
	}
	if err != nil {
		s.errCh <- err
		return nil
	}
	return pagers
}
