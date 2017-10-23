package postgres

import (
	"encoding/json"

	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
)

var _ watchub.StargazersSvc = &StargazerSvc{}

type StargazerSvc struct {
	db *sqlx.DB
}

func (s *StargazerSvc) Get(execution watchub.Execution) (result []watchub.Star, err error) {
	var stars json.RawMessage
	err = s.db.QueryRow(
		"SELECT stars FROM tokens WHERE user_id = $1",
		execution.UserID,
	).Scan(&stars)
	if err != nil {
		return result, err
	}
	return result, json.Unmarshal(stars, &result)
}
