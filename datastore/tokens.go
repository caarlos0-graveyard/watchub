package datastore

import "golang.org/x/oauth2"

type Tokenstore interface {
	SaveToken(userID int64, token *oauth2.Token) error
	Schedule(userID int64) error
}
