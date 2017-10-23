package postgres

import (
	"encoding/json"
	"time"

	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var _ watchub.TokensSvc = &TokensSvc{}

func NewTokensSvc(db *sqlx.DB) *TokensSvc {
	return &TokensSvc{
		db: db,
	}
}

type TokensSvc struct {
	db *sqlx.DB
}

func (s *TokensSvc) Schedule(userID int64, date time.Time) error {
	_, err := s.db.Exec(
		"UPDATE tokens SET next = $2, updated_at = now() WHERE user_id = $1",
		userID,
		date,
	)
	return err
}

func (s *TokensSvc) Exists(userID int64) (result bool, err error) {
	err = s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM tokens WHERE user_id = $1)",
		userID,
	).Scan(&result)
	return
}

const insertTokenStm = `
	INSERT INTO tokens(user_id, token)
	VALUES($1, $2)
	ON CONFLICT(user_id)
		DO UPDATE SET token = $2, updated_at = now();
`

func (s *TokensSvc) Save(userID int64, token *oauth2.Token) error {
	strToken, err := tokenToJSON(token)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(insertTokenStm, userID, strToken)
	return err
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	d, err := json.Marshal(token)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshall json token")
	}
	return string(d), nil
}
