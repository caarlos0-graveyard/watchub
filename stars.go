package watchub

// Star of a github repo
type Star struct {
	RepoID     int64    `json:"repo_id"`
	RepoName   string   `json:"repo_name"`
	Stargazers []string `json:"stargazers"`
}

type StargazersReadSvc interface {
	Get(execution Execution) ([]Star, error)
}

type StargazersCountSvc interface {
	Count(userID int64) (int, error)
}

type StargazersWriteSvc interface {
	Save(userID int64, stars []Star) error
}

// TODO: break this in smaller interfaces and compose
type StargazersSvc interface {
	StargazersWriteSvc
	StargazersReadSvc
	StargazersCountSvc
}
