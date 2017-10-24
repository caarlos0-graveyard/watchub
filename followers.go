package watchub

type FollowersSvc interface {
	Get(execution Execution) ([]string, error)
	Save(userID int64, followers []string) error
	Count(userID int64) (int, error)
}
