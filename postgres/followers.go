package postgres

import (
	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var _ watchub.FollowersSvc = &FollowersSvc{}

func NewFollowersSvc(db *sqlx.DB) *FollowersSvc {
	return &FollowersSvc{
		db: db,
	}
}

type FollowersSvc struct {
	db *sqlx.DB
}

func (s *FollowersSvc) Save(userID int64, followers watchub.Followers) error {
	_, err := s.db.Exec(
		"UPDATE tokens SET followers = $2 WHERE user_id = $1",
		userID,
		pq.Array(followers),
	)
	return err
}

func (s *FollowersSvc) Get(execution watchub.Execution) (watchub.Followers, error) {
	var logins []string
	return watchub.Followers(logins), s.db.QueryRow(
		"SELECT followers FROM tokens WHERE user_id = $1",
		execution.UserID,
	).Scan(pq.Array(&logins))
}

const followerCountQuery = `
	SELECT COALESCE(array_length(followers, 1), 0)
	FROM tokens
	WHERE user_id = $1
`

func (s *FollowersSvc) Count(userID int64) (count int, err error) {
	err = s.db.QueryRow(followerCountQuery, userID).Scan(&count)
	return
}
