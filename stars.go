package watchub

// Star of a github repo
type Star struct {
	RepoID     int64    `json:"repo_id"`
	RepoName   string   `json:"repo_name"`
	Stargazers []string `json:"stargazers"`
}

type StargazersSvc interface {
	Get(execution Execution) ([]Star, error)
}
