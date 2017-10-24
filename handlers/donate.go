package handlers

import (
	"net/http"

	"github.com/caarlos0/watchub"

	"github.com/caarlos0/watchub/config"
	"github.com/gorilla/sessions"
)

func NewDonate(
	config config.Config,
	session sessions.Store,
) *Donate {
	return &Donate{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

type Donate struct {
	Base
	stars        watchub.StargazersSvc
	followers    watchub.FollowersSvc
	repositories watchub.RepositoriesSvc
}

func (h *Donate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.render(w, "donate", h.sessionData(w, r))
}
