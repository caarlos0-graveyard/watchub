package handlers

import (
	"net/http"
	"time"

	"github.com/caarlos0/watchub"

	"github.com/caarlos0/watchub/config"
	"github.com/gorilla/sessions"
)

func NewSchedule(
	config config.Config,
	session sessions.Store,
	tokens watchub.TokensSvc,
) *Schedule {
	return &Schedule{
		Base: Base{
			config:  config,
			session: session,
		},
		tokens: tokens,
	}
}

type Schedule struct {
	Base
	tokens watchub.TokensSvc
}

func (h *Schedule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, _ := h.session.Get(r, h.config.SessionName)
	id, _ := session.Values["user_id"].(int)
	if session.IsNew || id == 0 {
		http.Error(w, "not logged in", http.StatusForbidden)
		return
	}
	if err := h.tokens.Schedule(int64(id), time.Now()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "scheduled", h.sessionData(w, r))
}
