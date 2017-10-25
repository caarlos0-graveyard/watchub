package watchub

type FollowersCountSvc interface {
	Count(userID int64) (int, error)
}

type FollowersReadSvc interface {
	Get(execution Execution) ([]string, error)
}

type FollowersWriteSvc interface {
	Save(userID int64, followers []string) error
}

// TODO: create a followers type here
type FollowersSvc interface {
	FollowersReadSvc
	FollowersCountSvc
	FollowersWriteSvc
}
