package database

import (
	"github.com/caarlos0/watchub/shared/model"
	"github.com/jmoiron/sqlx"
)

// Execstore in database
type Execstore struct {
	*sqlx.DB
}

// NewExecstore datastore
func NewExecstore(db *sqlx.DB) *Execstore {
	return &Execstore{db}
}

const executionsStmQuery = `
	UPDATE tokens
	SET
		next = datetime(current_timestamp, '+1 days'),
		updated_at = current_timestamp
	WHERE next <= current_timestamp and disabled is not true
	RETURNING id, token
`

// Executions get the executions that should be made
func (db *Execstore) Executions() (executions []model.Execution, err error) {
	return executions, db.Select(&executions, executionsStmQuery)
}
