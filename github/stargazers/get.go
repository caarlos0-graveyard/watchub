package stargazers

import (
	"context"
	"sync"
	"time"

	"github.com/caarlos0/watchub/shared/model"
	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Get the list of repos of a given user
func Get(
	ctx context.Context,
	client *github.Client,
	repos []*github.Repository,
) (result []model.Star, err error) {
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
) (result model.Star, err error) {
	stars, err := stars(ctx, client, repo)
	if err != nil {
		return result, err
	}
	// nolint: prealloc
	var stargazers []string
	for _, star := range stars {
		stargazers = append(stargazers, star.User.GetLogin())
	}
	return model.Star{
		RepoID:     *repo.ID,
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
	for {
		repos, nextPage, err := getPage(ctx, client, repo, opt)
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

func getPage(
	ctx context.Context,
	client *github.Client,
	repo *github.Repository,
	opt *github.ListOptions,
) ([]*github.Stargazer, int, error) {
	stars, resp, err := client.Activity.ListStargazers(
		ctx,
		*repo.Owner.Login,
		*repo.Name,
		opt,
	)
	var aerr github.AbuseRateLimitError
	if err != nil && errors.As(err, &aerr) {
		time.Sleep(aerr.GetRetryAfter())
		return getPage(ctx, client, repo, opt)
	}
	if err != nil {
		return stars, 0, err
	}
	return stars, resp.NextPage, nil
}
