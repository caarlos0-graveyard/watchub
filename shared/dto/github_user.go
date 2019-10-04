package dto

type GitHubUser struct {
	ID        int64
	Login     string
	Email     string
	Followers []string
}
