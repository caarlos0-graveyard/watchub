package model

// Star of a github repo
type Star struct {
	RepoName   string   `json:"repo_name"`
	Stargazers []string `json:"stargazers"`
}
