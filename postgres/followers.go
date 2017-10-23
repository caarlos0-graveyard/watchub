package postgres

import (
	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var _ watchub.FollowersSvc = &FollowersSvc{}

type FollowersSvc struct {
	db *sqlx.DB
}

func (s *FollowersSvc) Get(execution watchub.Execution) ([]string, error) {
	var logins []string
	return logins, s.db.QueryRow(
		"SELECT followers FROM tokens WHERE user_id = $1",
		execution.UserID,
	).Scan(pq.Array(&logins))
}
