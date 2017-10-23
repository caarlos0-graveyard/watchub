package handlers

import (
	"net/http"

	"github.com/caarlos0/watchub"

	"github.com/caarlos0/watchub/config"
	"github.com/gorilla/sessions"
)

// IndexPageData DTO
type IndexPageData struct {
	PageData
	FollowerCount   int
	StarCount       int
	RepositoryCount int
}

func NewIndex(
	config config.Config,
	session sessions.Store,
	stars watchub.StargazersSvc,
	followers watchub.FollowersSvc,
	repositories watchub.RepositoriesSvc,
) *Index {
	return &Index{
		Base: Base{
			config:  config,
			session: session,
		},
		stars:        stars,
		followers:    followers,
		repositories: repositories,
	}
}

type Index struct {
	Base
	stars        watchub.StargazersSvc
	followers    watchub.FollowersSvc
	repositories watchub.RepositoriesSvc
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data = IndexPageData{
		PageData: h.sessionData(w, r),
	}
	if data.User.ID > 0 {
		var err error
		var id = int64(data.User.ID)
		data.StarCount, err = h.stars.Count(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.FollowerCount, err = h.followers.Count(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.RepositoryCount, err = h.repositories.Count(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	h.render(w, "index", data)
}
