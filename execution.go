package watchub

// Execution model
type Execution struct {
	UserID int64  `db:"user_id"`
	Token  string `db:"token" json:"-"`
}

type ExecutionsSvc interface {
	All() ([]Execution, error)
}
