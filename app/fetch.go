package app

import (
	"context"
	"log"
	"star-golang-orms/domain/model"
	"star-golang-orms/domain/repository"
	"star-golang-orms/domain/service"
	"sync"
)

type fetchService struct {
	gitHubRepo  repository.GitHub
	wgR         sync.WaitGroup
	wgD         sync.WaitGroup
	mu          sync.Mutex
	repos       model.Repositories
	detailRepos model.RepositoryDetails
	errCh       chan error
}

func NewFetchService(repo repository.GitHub) service.Fetcher {
	return &fetchService{
		gitHubRepo:  repo,
		wgR:         sync.WaitGroup{},
		wgD:         sync.WaitGroup{},
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
			s.wgR.Add(1)
			go func(repoNm string) {
				defer s.wgR.Done()
				repo, err := s.gitHubRepo.GetRepository(ctx, model.RepositoryName(repoNm))
				if err != nil {
					s.errCh <- err
					return
				}
				s.repos = append(s.repos, *repo)

				log.Println(repoNm + " Start")
				stargazers := s.getStargazersCountByRepo(ctx, *repo)
				log.Println(repoNm + " DONE")

				s.mu.Lock()
				defer s.mu.Unlock()
				s.detailRepos = append(s.detailRepos, *model.NewRepositoryDetails(*repo, stargazers))
			}(repoNm)
		}
	}
	s.wgR.Wait()

	return nil
}

// func (s *fetchService) getStargazersCountByRepo(ctx context.Context, repo model.Repository) []model.Stargazer {
// 	var lock sync.Mutex
// 	var stargazers []model.Stargazer
// 	for page := 1; page <= model.LastPage(repo); page++ {
// 		s.wgD.Add(1)
// 		go func(page int) {
// 			defer s.wgD.Done()
// 			pagers, err := s.gitHubRepo.GetStarPage(ctx, repo, page)
// 			if errors.Is(err, domain.ErrNoMorePages) {
// 				return
// 			}
// 			if err != nil {
// 				s.errCh <- err
// 				return
// 			}
// 			lock.Lock()
// 			defer lock.Unlock()
// 			stargazers = append(stargazers, *pagers...)
// 		}(page)
// 	}
// 	s.wgD.Wait()

// 	return stargazers
// }

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

// gorutionをリファクタする
// appendをドメイン層でやる
// メソッド名、責務を考える
