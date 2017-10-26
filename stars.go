package watchub

// Star of a github repo
type Star struct {
	RepoID     int64    `json:"repo_id"`
	RepoName   string   `json:"repo_name"`
	Stargazers []string `json:"stargazers"`
}

type Stars []Star

func (stars Stars) Count() (count int) {
	for _, star := range stars {
		count += len(star.Stargazers)
	}
	return count
}

type StargazersReadSvc interface {
	Get(execution Execution) (Stars, error)
}

type StargazersCountSvc interface {
	Count(userID int64) (int, error)
}

type StargazersWriteSvc interface {
	Save(userID int64, stars Stars) error
}

type StargazersSvc interface {
	StargazersWriteSvc
	StargazersReadSvc
	StargazersCountSvc
}
