package database

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/datastore"
	"github.com/jmoiron/sqlx"
)

// Connect creates a connection pool to the database
func Connect() *sql.DB {
	url := "watchub.db"
	var log = log.WithField("url", url)
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		log.WithError(err).Fatal("Failed to open connection to database")
	}
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("Failed to ping database")
	}
	return db
}

// NewDatastore returns a new Datastore
func NewDatastore(db *sql.DB) datastore.Datastore {
	var dbx = sqlx.NewDb(db, "sqlite3")
	return struct {
		*Tokenstore
		*Execstore
		*Userdatastore
	}{
		NewTokenstore(dbx),
		NewExecstore(dbx),
		NewUserdatastore(dbx),
	}
}
