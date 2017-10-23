package watchub

type FollowersSvc interface {
	Get(execution Execution) ([]string, error)
	Count(userID int64) (int, error)
}
