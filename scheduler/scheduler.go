package scheduler

import (
	"time"

	"github.com/apex/log"
	"github.com/caarlos0/watchub"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/diff"
	"github.com/caarlos0/watchub/oauth"
	"github.com/gorilla/sessions"
	"github.com/robfig/cron"
)

var _ watchub.ScheduleSvc = &Scheduler{}

// Scheduler type
type Scheduler struct {
	cron              *cron.Cron
	Config            config.Config
	Oauth             *oauth.Oauth
	Session           sessions.Store
	Mailer            watchub.MailSvc
	Users             watchub.UsersSvc
	Executions        watchub.ExecutionsSvc
	PreviousStars     watchub.StargazersSvc
	PreviousFollowers watchub.FollowersSvc
	CurrentStars      watchub.StargazersSvc
	CurrentFollowers  watchub.FollowersSvc
}

// Start the scheduler
func (s *Scheduler) Start() {
	var fn = func() {
		execs, err := s.Executions.All()
		if err != nil {
			log.WithError(err).Error("failed to get executions")
			return
		}
		for _, exec := range execs {
			exec := exec
			go s.process(exec)
		}
	}
	s.cron = cron.New()
	if err := s.cron.AddFunc(s.Config.Schedule, fn); err != nil {
		log.WithError(err).Fatal("failed to start cron service")
	}
	s.cron.Start()
}

// Stop the scheduler
func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
	}
	s.cron = nil
}

func (s *Scheduler) process(exec watchub.Execution) {
	var start = time.Now()
	defer log.WithField("time_taken", time.Since(start).Seconds()).Info("done")
	var log = log.WithField("id", exec.UserID)
	previousStars, err := s.PreviousStars.Get(exec)
	if err != nil {
		log.WithError(err).Error("failed")
		return
	}
	previousFollowers, err := s.PreviousFollowers.Get(exec)
	if err != nil {
		log.WithError(err).Error("failed")
		return
	}
	currentStars, err := s.CurrentStars.Get(exec)
	if err != nil {
		log.WithError(err).Error("failed")
		return
	}
	currentFollowers, err := s.CurrentFollowers.Get(exec)
	if err != nil {
		log.WithError(err).Error("failed")
		return
	}
	if err := s.CurrentStars.Save(exec.UserID, currentStars); err != nil {
		log.WithError(err).Error("failed")
		return
	}
	if err := s.CurrentFollowers.Save(exec.UserID, currentFollowers); err != nil {
		log.WithError(err).Error("failed")
		return
	}
	user, err := s.Users.Info(exec.Token)
	if err != nil {
		log.WithError(err).Error("failed")
		return
	}

	// new user, welcome him!
	if len(previousFollowers)+len(previousStars) == 0 {
		s.Mailer.SendWelcome(watchub.WelcomeEmail{
			Login:     user.Login,
			Email:     user.Email,
			Followers: len(currentFollowers),
			Stars:     countStars(currentStars),
			Repos:     len(currentStars),
		})
		return
	}
	newFollowers := diff.Of(currentFollowers, previousFollowers)
	unfollowers := diff.Of(previousFollowers, currentFollowers)
	newStars, unstars := stargazerStatistics(currentStars, previousStars)
	// nothing changed, ignore
	if len(newFollowers)+len(unfollowers)+len(newStars)+len(unstars) == 0 {
		return
	}

	// send changes!
	s.Mailer.SendChanges(
		watchub.ChangesEmail{
			Login:        user.Login,
			Email:        user.Email,
			Followers:    len(currentFollowers),
			Stars:        countStars(currentStars),
			Repos:        len(currentStars),
			NewFollowers: newFollowers,
			Unfollowers:  unfollowers,
			NewStars:     newStars,
			Unstars:      unstars,
		},
	)
}

// TODO: refactor this to a watchub.Stars type?
func countStars(stars []watchub.Star) (count int) {
	for _, star := range stars {
		count += len(star.Stargazers)
	}
	return
}

func stargazerStatistics(stars, previousStars []watchub.Star) (newStars, unstars []watchub.StarEmail) {
	for _, s := range stars {
		for _, os := range previousStars {
			if s.RepoID != os.RepoID {
				continue
			}
			if d := getDiff(s.RepoName, s.Stargazers, os.Stargazers); d != nil {
				newStars = append(newStars, *d)
			}
			if d := getDiff(s.RepoName, os.Stargazers, s.Stargazers); d != nil {
				unstars = append(unstars, *d)
			}
			break
		}
	}
	return newStars, unstars
}

func getDiff(name string, x, xs []string) *watchub.StarEmail {
	var d = diff.Of(x, xs)
	if len(d) > 0 {
		return &watchub.StarEmail{
			Repo:  name,
			Users: d,
		}
	}
	return nil
}
