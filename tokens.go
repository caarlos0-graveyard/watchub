package watchub

import (
	"time"

	"golang.org/x/oauth2"
)

type TokensSvc interface {
	Exists(userID int64) (bool, error)
	Save(userID int64, token *oauth2.Token) error
	Schedule(userID int64, date time.Time) error
}
