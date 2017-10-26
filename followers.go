package watchub

type FollowersCountSvc interface {
	Count(userID int64) (int, error)
}

type FollowersReadSvc interface {
	Get(execution Execution) (Followers, error)
}

type FollowersWriteSvc interface {
	Save(userID int64, followers Followers) error
}

type Followers []string

type FollowersSvc interface {
	FollowersReadSvc
	FollowersCountSvc
	FollowersWriteSvc
}
