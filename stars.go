package watchub

// Star of a github repo
type Star struct {
	RepoID     int64    `json:"repo_id"`
	RepoName   string   `json:"repo_name"`
	Stargazers []string `json:"stargazers"`
}

// TODO: break this in smaller interfaces and compose
type StargazersSvc interface {
	Get(execution Execution) ([]Star, error)
	Count(userID int64) (int, error)
	Save(userID int64, stars []Star) error
}
