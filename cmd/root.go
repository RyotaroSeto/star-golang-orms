package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"star-golang-orms/configs"
	"star-golang-orms/pkg"
	"sync"
	"syscall"
)

func Execute() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
		return
	}

	gh, err := ExecGitHubAPI(config.GithubToken)
	if err != nil {
		log.Fatal("cannot exec github api", err)
		return
	}

	err = gh.SortDesByStarCount()
	if err != nil {
		log.Fatal("cannot sort star count", err)
		return
	}

	err = gh.MakeChart()
	if err != nil {
		log.Fatal("connot make chart", err)
		return
	}

	err = gh.Edit()
	if err != nil {
		log.Fatal("connot edit readme", err)
		return
	}
}

func ExecGitHubAPI(token string) (pkg.GitHub, error) {
	ctx, cancel := NewCtx()
	defer cancel()

	var repos []pkg.GithubRepository
	var detaiRepos []pkg.ReadmeDetailsRepository

	wg := new(sync.WaitGroup)
	var lock sync.Mutex
	for _, repoNm := range pkg.TargetRepository {
		wg.Add(1)
		fmt.Println("start:" + repoNm)
		go func(repoNm string) {
			defer wg.Done()
			repo, err := pkg.NowGithubRepoCount(ctx, repoNm, token)
			if err != nil {
				log.Println(err)
				return
			}
			repos = append(repos, repo)
			fmt.Println(repoNm + " Start")
			repo, stargazers := pkg.GetRepo(ctx, repoNm, token, repo)
			fmt.Println(repoNm + " DONE")
			lock.Lock()
			defer lock.Unlock()
			detaiRepos = append(detaiRepos, pkg.NewDetailsRepository(repo, stargazers))
		}(repoNm)
	}

	wg.Wait()
	gh := pkg.NewGitHub(repos, detaiRepos)
	return gh, nil
}

func NewCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
		<-trap
	}()

	return ctx, cancel
}
