package watchub

// TODO: create a followers type here
// TODO: break this in smaller interfaces and compose
type FollowersSvc interface {
	Get(execution Execution) ([]string, error)
	Save(userID int64, followers []string) error
	Count(userID int64) (int, error)
}
