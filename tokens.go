package watchub

import (
	"time"

	"golang.org/x/oauth2"
)

// TODO: break this in smaller interfaces and compose
type TokensSvc interface {
	Exists(userID int64) (bool, error)
	Save(userID int64, token *oauth2.Token) error
	Schedule(userID int64, date time.Time) error
}
