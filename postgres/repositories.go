package postgres

import (
	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
)

var _ watchub.RepositoriesSvc = &RepositoriesSvc{}

func NewRepositoriesSvc(db *sqlx.DB) *RepositoriesSvc {
	return &RepositoriesSvc{
		db: db,
	}
}

type RepositoriesSvc struct {
	db *sqlx.DB
}

var repositoryCountQuery = `
	SELECT COALESCE(json_array_length(stars), 0)
	FROM tokens
	WHERE user_id = $1
`

func (s *RepositoriesSvc) Count(userID int64) (count int, err error) {
	err = s.db.QueryRow(repositoryCountQuery, userID).Scan(&count)
	return
}
