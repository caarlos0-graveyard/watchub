package watchub

// WelcomeEmail is the DTO passed to the welcome email template
type WelcomeEmail struct {
	Login     string
	Email     string
	Followers int
	Stars     int
	Repos     int
}

// StarEmail is the DTO representing a repository and the users starring it
type StarEmail struct {
	Repo  string
	Users []string
}

// ChangesEmail is the DTO passed down to the daily email
type ChangesEmail struct {
	Login        string
	Email        string
	Followers    int
	Stars        int
	Repos        int
	NewFollowers []string
	Unfollowers  []string
	NewStars     []StarEmail
	Unstars      []StarEmail
}

type MailSvc interface {
	SendWelcome(data WelcomeEmail)
	SendChanges(data ChangesEmail)
}
