package handlers

import (
	"net/http"

	"github.com/caarlos0/watchub/oauth"
)

func NewLogin(oauth *oauth.Oauth) *Login {
	return &Login{oauth}
}

type Login struct {
	oauth *oauth.Oauth
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var url = h.oauth.AuthCodeURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
