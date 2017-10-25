package github

import (
	"context"

	"github.com/caarlos0/watchub"
	"github.com/caarlos0/watchub/oauth"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

var _ watchub.FollowersReadSvc = &FollowersSvc{}

func NewFollowersSvc(oauth *oauth.Oauth) *FollowersSvc {
	return &FollowersSvc{
		oauth: oauth,
	}
}

type FollowersSvc struct {
	oauth *oauth.Oauth
}

func (s *FollowersSvc) Get(execution watchub.Execution) ([]string, error) {
	var result []string
	var ctx = context.Background()
	client, err := s.oauth.ClientFrom(ctx, execution.Token)
	if err != nil {
		return result, err
	}
	var getPage = func(opt *github.ListOptions) (followers []*github.User, nextPage int, err error) {
		followers, resp, err := client.Users.ListFollowers(ctx, "", opt)
		if err != nil {
			return
		}
		return followers, resp.NextPage, err
	}
	var opt = &github.ListOptions{PerPage: 30}
	for {
		followers, nextPage, err := getPage(opt)
		if err != nil {
			return result, errors.Wrap(err, "failed to get followers")
		}
		for _, u := range followers {
			result = append(result, u.GetLogin())
		}
		if opt.Page = nextPage; nextPage == 0 {
			break
		}
	}
	return result, nil
}
