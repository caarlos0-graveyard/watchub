package model

// Execution model
type Execution struct {
	UserID int64  `db:"id"`
	Token  string `db:"token" json:"-"`
}
