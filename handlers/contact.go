package handlers

import (
	"net/http"

	"github.com/caarlos0/watchub"

	"github.com/caarlos0/watchub/config"
	"github.com/gorilla/sessions"
)

func NewContact(
	config config.Config,
	session sessions.Store,
) *Contact {
	return &Contact{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

type Contact struct {
	Base
	stars        watchub.StargazersSvc
	followers    watchub.FollowersSvc
	repositories watchub.RepositoriesSvc
}

func (h *Contact) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render(w, "contact", h.sessionData(w, r))
}
