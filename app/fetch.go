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

	gh.ReadmeRepoAndDetailSort()
	if err = gh.MakeHTMLChartFile(); err != nil {
		return starError.Newf(starError.InternalServerError, "failed to create html file: %s", err)
	}

	if err = model.ConvertHTMLToImage(); err != nil {
		return starError.Newf(starError.InternalServerError, "failed to convert html to image: %s", err)
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
	wg := new(sync.WaitGroup)
	var lock sync.Mutex
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-s.errCh:
		return err
	default:
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
	}
	wg.Wait()

	return nil

}

func (s *fetchService) getStargazersCountByRepo(ctx context.Context, repo *model.Repository) []model.Stargazer {
	var (
		sem        = make(chan bool, 4)
		eg         errgroup.Group
		lock       sync.Mutex
		stargazers = model.NewStargazers()
	)
	for page := 1; page <= model.LastPage(repo); page++ {
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result := s.fetchStargazersPage(ctx, repo, page, stargazers)
			lock.Lock()
			defer lock.Unlock()
			stargazers.Stars = append(stargazers.Stars, *result...)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		log.Println(err)
	}

	return stargazers.Stars
}

func (s *fetchService) fetchStargazersPage(ctx context.Context, repo *model.Repository, page int, stargazers *model.Stargazers) *[]model.Stargazer {
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

// READMEの結果がソートされていない
