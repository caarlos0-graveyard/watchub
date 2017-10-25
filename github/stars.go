package github

import (
	"context"
	"sync"

	"github.com/caarlos0/watchub"
	"github.com/caarlos0/watchub/oauth"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var _ watchub.StargazersReadSvc = &StargazersSvc{}

func NewStargazersSvc(oauth *oauth.Oauth) *StargazersSvc {
	return &StargazersSvc{
		oauth: oauth,
	}
}

type StargazersSvc struct {
	oauth *oauth.Oauth
}

func (s *StargazersSvc) Get(execution watchub.Execution) (result []watchub.Star, err error) {
	var ctx = context.Background()
	client, err := s.oauth.ClientFrom(ctx, execution.Token)
	if err != nil {
		return result, err
	}
	repos, err := getRepos(ctx, client)
	if err != nil {
		return result, err
	}
	return getStars(ctx, client, repos)
}

func getRepos(
	ctx context.Context,
	client *github.Client,
) (result []*github.Repository, err error) {
	var opt = &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var getPage = func(opt *github.RepositoryListOptions) (repos []*github.Repository, nextPage int, err error) {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return
		}
		return repos, resp.NextPage, err
	}
	for {
		repos, nextPage, err := getPage(opt)
		if err != nil {
			return result, errors.Wrap(err, "failed to get repositories")
		}
		for _, repo := range repos {
			if repo.GetFork() || repo.GetPrivate() {
				continue
			}
			result = append(result, repo)
		}
		if opt.Page = nextPage; nextPage == 0 {
			break
		}
	}
	return
}

func getStars(
	ctx context.Context,
	client *github.Client,
	repos []*github.Repository,
) (result []watchub.Star, err error) {
	var g errgroup.Group
	var m sync.Mutex
	var pool = make(chan bool, 5)
	for _, repo := range repos {
		repo := repo
		pool <- true
		g.Go(func() error {
			defer func() {
				<-pool
			}()
			r, er := processRepo(ctx, client, repo)
			if er != nil {
				return errors.Wrap(er, "failed to get repository stars")
			}
			m.Lock()
			defer m.Unlock()
			result = append(result, r)
			return nil
		})
	}
	err = g.Wait()
	return
}

func processRepo(
	ctx context.Context,
	client *github.Client,
	repo *github.Repository,
) (result watchub.Star, err error) {
	stars, err := stars(ctx, client, repo)
	if err != nil {
		return result, err
	}
	var stargazers []string
	for _, star := range stars {
		stargazers = append(stargazers, *star.User.Login)
	}
	return watchub.Star{
		RepoID:     int64(*repo.ID),
		RepoName:   *repo.FullName,
		Stargazers: stargazers,
	}, nil
}

func stars(
	ctx context.Context,
	client *github.Client,
	repo *github.Repository,
) (result []*github.Stargazer, err error) {
	var opt = &github.ListOptions{PerPage: 30}
	var getPage = func(opt *github.ListOptions) (stars []*github.Stargazer, nextPage int, err error) {
		stars, resp, err := client.Activity.ListStargazers(
			ctx, *repo.Owner.Login, *repo.Name, opt,
		)
		if err != nil {
			return
		}
		return stars, resp.NextPage, nil
	}
	for {
		repos, nextPage, err := getPage(opt)
		if err != nil {
			return result, err
		}
		result = append(result, repos...)
		if opt.Page = nextPage; nextPage == 0 {
			break
		}
	}
	return result, nil
}
