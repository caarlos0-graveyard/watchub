package postgres

import (
	"encoding/json"

	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
)

var _ watchub.StargazersSvc = &StargazersSvc{}

func NewStargazersSvc(db *sqlx.DB) *StargazersSvc {
	return &StargazersSvc{
		db: db,
	}
}

type StargazersSvc struct {
	db *sqlx.DB
}

func (s *StargazersSvc) Get(execution watchub.Execution) (result watchub.Stars, err error) {
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

const starCountQuery = `
	SELECT COALESCE(SUM(json_array_length((repo->>'stargazers')::json)), 0)
	FROM tokens t
	CROSS JOIN json_array_elements(t.stars) repo
	WHERE t.user_id = $1
`

func (s *StargazersSvc) Count(userID int64) (count int, err error) {
	err = s.db.QueryRow(starCountQuery, userID).Scan(&count)
	return
}

func (s *StargazersSvc) Save(userID int64, stars watchub.Stars) error {
	data, err := json.Marshal(stars)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		"UPDATE tokens SET stars = $2 WHERE user_id = $1",
		userID,
		data,
	)
	return err
}
