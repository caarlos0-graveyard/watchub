package postgres

import (
	"github.com/caarlos0/watchub"
	"github.com/jmoiron/sqlx"
)

var _ watchub.ExecutionsSvc = &ExecutionsSvc{}

func NewExecutionsSvc(db *sqlx.DB) *ExecutionsSvc {
	return &ExecutionsSvc{
		db: db,
	}
}

type ExecutionsSvc struct {
	db *sqlx.DB
}

const executionsStmQuery = `
	UPDATE tokens
	SET next = now() + interval '1 day', updated_at = now()
	WHERE next <= now()
	RETURNING user_id, token
`

func (s *ExecutionsSvc) All() (executions []watchub.Execution, err error) {
	return executions, s.db.Select(&executions, executionsStmQuery)
}
