package repos

import (
	"context"

	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
)

// Get all user's repos
func Get(
	ctx context.Context,
	client *github.Client,
) (result []*github.Repository, err error) {
	var opt = &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
	for {
		repos, nextPage, err := getPage(ctx, client, opt)
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

func getPage(
	ctx context.Context,
	client *github.Client,
	opt *github.RepositoryListOptions,
) (repos []*github.Repository, nextPage int, err error) {
	repos, resp, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		return
	}
	return repos, resp.NextPage, err
}
