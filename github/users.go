package github

import (
	"context"

	"github.com/caarlos0/watchub"
	"github.com/caarlos0/watchub/oauth"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

var _ watchub.UsersSvc = &Usersvc{}

func NewUsersSvc(oauth *oauth.Oauth) *UsersSvc {
	return &UsersSvc{
		oauth: oauth,
	}
}

type UsersSvc struct {
	oauth *oauth.Oauth
}

func (s *UsersSvc) Info(token string) (watchub.UserInfo, error) {
	var info watchub.UserInfo
	var ctx = context.Background()
	client, err := s.oauth.ClientFrom(ctx, token)
	if err != nil {
		return info, err
	}
	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return info, err
	}
	info.Login = u.GetLogin()
	emails, _, err := client.Users.ListEmails(ctx, &github.ListOptions{PerPage: 10})
	if err != nil {
		return info, errors.Wrap(err, "failed to get current user email addresses")
	}
	for _, e := range emails {
		if e.GetPrimary() && e.GetVerified() {
			info.Email = e.GetEmail()
			break
		}
	}
	return info, nil
}
