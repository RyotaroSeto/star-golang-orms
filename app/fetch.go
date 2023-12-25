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

	gh.ReadmeRepoAndDetailSort()
	if err = gh.MakeHTMLChartFile(); err != nil {
		return err
	}

	if err = model.ConvertHTMLToImage(); err != nil {
		return err
	}

	return gh.ReadmeEdit()
}

func (s *fetchService) ExecGitHubAPI(ctx context.Context, token string) (*model.GitHub, error) {
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
			stargazers := GetStargazersCountByRepo(ctx, token, repo)
			// stargazers := s.GetStargazersCountByRepo(ctx, token, repo)
			log.Println(repoNm + " DONE")

			rd := model.NewRepositoryDetails(repo, stargazers)
			lock.Lock()
			defer lock.Unlock()
			s.detailRepos = append(s.detailRepos, rd)
		}(repoNm)
	}

	wg.Wait()
	// gh := pkg.NewGitHub(s.repos, detaiRepos)
	return &model.GitHub{
		Repositories:      s.repos,
		RepositoryDetails: s.detailRepos,
	}, nil
}

func GetStargazersCountByRepo(ctx context.Context, token string, repo *model.Repository) []model.Stargazer {
	sem := make(chan bool, 4)
	var eg errgroup.Group
	var lock sync.Mutex
	var stargazers []model.Stargazer
	for page := 1; page <= pkg.LastPage(repo); page++ {
		sem <- true
		page := page
		eg.Go(func() error {
			defer func() { <-sem }()
			result, err := pkg.GetStargazersPage(ctx, repo, page, token)
			// if errors.Is(err, ErrNoMorePages) {
			// 	return nil
			// }
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			stargazers = append(stargazers, result...)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		log.Println(err)
	}

	return stargazers
}
