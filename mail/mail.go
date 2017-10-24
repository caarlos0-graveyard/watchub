package mail

import (
	"fmt"

	"github.com/caarlos0/watchub"
	"github.com/matcornic/hermes"
)

var _ watchub.MailSvc = &MailSvc{}

func NewMailSvc() *MailSvc {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Watchub",
			Link: "http://watchub.pw",
			Logo: "https://raw.githubusercontent.com/caarlos0/watchub/master/static/apple-touch-icon-144x144.png",
		},
	}
	return &MailSvc{
		h: h,
	}
}

type MailSvc struct {
	h hermes.Hermes
}

const welcomeTempl = `
<p>
	You have <strong>%d</strong> followers and
	<strong>%d</strong> stars across <strong>%d</strong> repositories.
</p>
`

func (s *MailSvc) SendWelcome(data watchub.WelcomeEmail) {
	var md = fmt.Sprintf(welcomeTempl, data.Followers, data.Stars, data.Repos)
	var email = hermes.Email{
		Body: hermes.Body{
			Name: data.Login,
			Intros: []string{
				"Welcome to Watchub! We're very excited to have you on board.",
			},
			FreeMarkdown: hermes.Markdown(md),
		},
	}
}
func (s *MailSvc) SendChanges(data watchub.ChangesEmail) {

}
