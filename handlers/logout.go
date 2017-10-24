package handlers

import (
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/gorilla/sessions"
)

func NewLogout(config config.Config, session sessions.Store) *Logout {
	return &Logout{
		config:  config,
		session: session,
	}
}

type Logout struct {
	config  config.Config
	session sessions.Store
}

func (h *Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, _ := h.session.Get(r, h.config.SessionName)
	session.Values = map[interface{}]interface{}{}
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "", http.StatusTemporaryRedirect)
}
