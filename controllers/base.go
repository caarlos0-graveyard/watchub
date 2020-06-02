package controllers

import (
	"net/http"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/gorilla/sessions"
)

type Base struct {
	session sessions.Store
	config  config.Config
}

func (ctrl Base) sessionData(w http.ResponseWriter, r *http.Request) dto.PageData {
	var user dto.PageUserData
	session, err := ctrl.session.Get(r, ctrl.config.SessionName)
	if err != nil {
		log.WithError(err).Errorf("failed to get session")
	}

	log.WithFields(log.Fields{
		"user_id":    session.Values["user_id"],
		"user_login": session.Values["user_login"],
		"new_user":   session.Values["new_user"],
		"isNew":      session.IsNew,
	}).Info("got session")

	user.ID, _ = session.Values["user_id"].(int64)
	user.Login, _ = session.Values["user_login"].(string)
	user.IsNew, _ = session.Values["new_user"].(bool)
	delete(session.Values, "new_user")

	if err := session.Save(r, w); err != nil {
		log.WithError(err).Error("failed to update session")
	}

	log.WithFields(log.Fields{
		"user_id":    user.ID,
		"user_login": user.Login,
		"new_user":   user.IsNew,
	}).Info("dto")

	return dto.PageData{
		User:     user,
		ClientID: ctrl.config.ClientID,
	}
}
