package database

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Tokenstore in database
type Tokenstore struct {
	*sqlx.DB
}

// NewTokenstore datastore
func NewTokenstore(db *sqlx.DB) *Tokenstore {
	return &Tokenstore{db}
}

const insertTokenStm = `
	INSERT INTO tokens(user_id, token)
	VALUES($1, $2)
	ON CONFLICT(user_id)
		DO UPDATE SET token = $2, updated_at = now(), disabled = false;
`

// SaveToken saves a token
func (db *Tokenstore) SaveToken(userID int64, token *oauth2.Token) error {
	strToken, err := tokenToJSON(token)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertTokenStm, userID, strToken)
	return err
}

// Schedule schedules a new execution at the given time
func (db *Tokenstore) Schedule(userID int64, date time.Time) error {
	_, err := db.Exec(
		"UPDATE tokens SET next = $2, updated_at = now() WHERE user_id = $1",
		userID,
		date,
	)
	return err
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	d, err := json.Marshal(token)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshall json token")
	}
	return string(d), nil
}
