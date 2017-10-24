package mail

import (
	"html/template"
	"io/ioutil"
	"testing"

	"github.com/caarlos0/watchub"
	"github.com/caarlos0/watchub/config"
	"github.com/tj/assert"
)

func TestWelcomeMail(t *testing.T) {
	s := MailSvc{
		hermes: emailConfig,
		config: config.Config{
			ClientID: "1",
		},
		welcome: template.Must(template.ParseFiles("../static/mail/welcome.md")),
	}
	data := watchub.WelcomeEmail{
		Email:     "caarlos0@gmail.com",
		Followers: 10,
		Login:     "caarlos0",
		Repos:     5,
		Stars:     1,
	}
	html, err := s.generate(data.Login, data, s.welcome, welcomeIntro)
	assert.NoError(t, err)
	// TODO: make this flag-toggleable
	// assert.NoError(t, ioutil.WriteFile("testdata/welcome.html", []byte(html), 0644))
	bts, err := ioutil.ReadFile("testdata/welcome.html")
	assert.NoError(t, err)
	assert.Equal(t, string(bts), html)
}

func TestChangesMail(t *testing.T) {
	s := MailSvc{
		hermes: emailConfig,
		config: config.Config{
			ClientID: "1",
		},
		changes: template.Must(template.ParseFiles("../static/mail/changes.md")),
	}
	data := watchub.ChangesEmail{
		Email:        "caarlos0@gmail.com",
		Followers:    10,
		Login:        "caarlos0",
		Repos:        5,
		Stars:        1,
		NewFollowers: []string{"juvenal", "moises"},
		NewStars: []watchub.StarEmail{
			{
				Repo:  "test/test",
				Users: []string{"juvenal", "moises"},
			},
		},
		Unfollowers: []string{"outro-juvenal", "outro-moises"},
		Unstars: []watchub.StarEmail{
			{
				Repo:  "test/test",
				Users: []string{"outro-juvenal", "outro-moises"},
			},
		},
	}
	html, err := s.generate(data.Login, data, s.changes, changesIntro)
	assert.NoError(t, err)
	// TODO: make this flag-toggleable
	// assert.NoError(t, ioutil.WriteFile("testdata/changes.html", []byte(html), 0644))
	bts, err := ioutil.ReadFile("testdata/changes.html")
	assert.NoError(t, err)
	assert.Equal(t, string(bts), html)
}
