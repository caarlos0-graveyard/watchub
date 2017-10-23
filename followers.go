package watchub

type FollowersSvc interface {
	Get(execution Execution) ([]string, error)
}
