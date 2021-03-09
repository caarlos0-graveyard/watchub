package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/caarlos0/watchub/github/email"
	"github.com/caarlos0/watchub/github/followers"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/google/go-github/v28/github"
)

var ErrBadCreds = errors.New("bad credentials")

// Info gets a github user info, like login, email and followers
func Info(ctx context.Context, client *github.Client) (user dto.GitHubUser, err error) {
	u, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		if resp.StatusCode == 401 {
			return user, fmt.Errorf("%w: %s", ErrBadCreds, err.Error())
		}
		return user, err
	}
	email, err := email.Get(ctx, client)
	if err != nil {
		return user, err
	}
	followers, err := followers.Get(ctx, client)
	if err != nil {
		return user, err
	}

	user.ID = u.GetID()
	user.Login = u.GetLogin()
	user.Email = email
	user.Followers = ToLoginArray(followers)
	return
}
