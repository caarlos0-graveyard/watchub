package oauth

import (
	"context"
	"encoding/json"

	"github.com/caarlos0/watchub/config"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Oauth info
type Oauth struct {
	config *oauth2.Config
	state  string
}

// New oauth
func New(
	config config.Config,
) *Oauth {
	return &Oauth{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{"user:email,public_repo"},
			Endpoint:     githuboauth.Endpoint,
		},
		state: config.OauthState,
	}
}

// Client for a given token
func (o *Oauth) Client(ctx context.Context, token *oauth2.Token) *github.Client {
	return github.NewClient(o.config.Client(ctx, token))
}

// ClientFrom for a given string token
func (o *Oauth) ClientFrom(ctx context.Context, tokenStr string) (*github.Client, error) {
	token, err := TokenFromJSON(tokenStr)
	if err != nil {
		return nil, err
	}
	return o.Client(ctx, token), nil
}

// AuthCodeURL URL to OAuth 2.0 provider's consent page
func (o *Oauth) AuthCodeURL() string {
	return o.config.AuthCodeURL(o.state, oauth2.AccessTypeOnline)
}

// IsStateValid true if state is valid
func (o *Oauth) IsStateValid(state string) bool {
	return o.state == state
}

// Exchange oauth code
func (o *Oauth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.config.Exchange(ctx, code)
}

// TokenFromJSON extract an oauth from a json string
func TokenFromJSON(str string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(str), &token); err != nil {
		return nil, errors.Wrap(err, "failed unmarshall json token")
	}
	return &token, nil
}
