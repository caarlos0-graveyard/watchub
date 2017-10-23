package watchub

type RepositoriesSvc interface {
	Count(userID int64) (int, error)
}
