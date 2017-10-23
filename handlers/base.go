package handlers

import (
	"html/template"
	"net/http"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/config"
	"github.com/gorilla/sessions"
)

type PageData struct {
	User     PageUserData
	ClientID string
}

type PageUserData struct {
	ID    int
	Login string
	IsNew bool
}

type Base struct {
	session sessions.Store
	config  config.Config
}

func (h Base) sessionData(w http.ResponseWriter, r *http.Request) PageData {
	var user PageUserData
	session, _ := h.session.Get(r, h.config.SessionName)
	if !session.IsNew {
		user.ID, _ = session.Values["user_id"].(int)
		user.Login, _ = session.Values["user_login"].(string)
		user.IsNew, _ = session.Values["new_user"].(bool)
		delete(session.Values, "new_user")
		if err := session.Save(r, w); err != nil {
			log.WithError(err).Error("failed to update session")
		}
	}
	return PageData{
		User:     user,
		ClientID: h.config.ClientID,
	}
}

func (h Base) render(w http.ResponseWriter, name string, data interface{}) {
	var t = template.Must(template.ParseFiles("static/layout.html", "static/"+name+".html"))
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
