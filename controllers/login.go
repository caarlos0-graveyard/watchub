package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/caarlos0/watchub/oauth"
	"github.com/gorilla/sessions"
)

// Login ctrl
type Login struct {
	Base
	oauth *oauth.Oauth
	store datastore.Datastore
}

// NewLogin ctrl
func NewLogin(
	config config.Config,
	session sessions.Store,
	oauth *oauth.Oauth,
	store datastore.Datastore,
) *Login {
	return &Login{
		Base: Base{
			config:  config,
			session: session,
		},
		store: store,
		oauth: oauth,
	}
}

// Handler handles /
func (ctrl *Login) Handler(w http.ResponseWriter, r *http.Request) {
	var url = ctrl.oauth.AuthCodeURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// CallbackHandler handles /login/callback
func (ctrl *Login) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	var state = r.FormValue("state")
	var code = r.FormValue("code")
	var ctx = context.Background()
	if !ctrl.oauth.IsStateValid(state) {
		log.Error("oauth state is invalid")
		http.Error(w, "invalid oauth state", http.StatusUnauthorized)
		return
	}
	token, err := ctrl.oauth.Exchange(ctx, code)
	if err != nil {
		log.WithError(err).Error("failed to exchange oauth token")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var client = ctrl.oauth.Client(ctx, token)
	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.WithError(err).Error("failed to authorize")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	exists, _ := ctrl.store.UserExist(int64(*u.ID))
	if err := ctrl.store.SaveToken(int64(*u.ID), token); err != nil {
		log.WithError(err).Error("failed to save token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		if err := ctrl.store.Schedule(int64(*u.ID), time.Now()); err != nil {
			log.WithError(err).Error("failed to schedule sync")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	session, err := ctrl.session.Get(r, ctrl.config.SessionName)
	if err != nil {
		log.WithError(err).Error("fail to get session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = *u.ID
	session.Values["user_login"] = *u.Login
	session.Values["new_user"] = !exists
	if err := session.Save(r, w); err != nil {
		log.WithError(err).Error("fail to save session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.WithFields(log.Fields{
		"user_id":    session.Values["user_id"],
		"user_login": session.Values["user_login"],
		"new_user":   session.Values["new_user"],
	}).Info("login successful")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
