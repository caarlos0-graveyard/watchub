package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/caarlos0/watchub"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/oauth"
)

func NewLoginCallback(
	oauth *oauth.Oauth,
	tokens watchub.TokensSvc,
	session sessions.Store,
	config config.Config,
) *LoginCallback {
	return &LoginCallback{
		oauth:   oauth,
		tokens:  tokens,
		session: session,
		config:  config,
	}
}

type LoginCallback struct {
	oauth   *oauth.Oauth
	tokens  watchub.TokensSvc
	session sessions.Store
	config  config.Config
}

func (h *LoginCallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var state = r.FormValue("state")
	var code = r.FormValue("code")
	var ctx = context.Background()
	if !h.oauth.IsStateValid(state) {
		http.Error(w, "invalid oauth state", http.StatusUnauthorized)
		return
	}
	token, err := h.oauth.Exchange(ctx, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var client = h.oauth.Client(ctx, token)
	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var userID = int64(u.GetID())
	exists, err := h.tokens.Exists(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.tokens.Save(userID, token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		if err := h.tokens.Schedule(userID, time.Now()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	session, _ := h.session.Get(r, h.config.SessionName)
	session.Values["user_id"] = int(userID)
	session.Values["user_login"] = u.GetLogin()
	session.Values["new_user"] = !exists
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
