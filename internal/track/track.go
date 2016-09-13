package track

import (
	"github.com/caarlos0/watchub/internal/config"
	mixpanel "github.com/nitrous-io/go-mixpanel"
)

func New(c config.Config) *Tracker {
	return &Tracker{
		mc: mixpanel.NewMixpanelClient(c.MixpanelToken),
	}
}

type Tracker struct {
	mc *mixpanel.Mixpanel
}

func (t *Tracker) Track(id int64, action string) error {
	return t.mc.Track(action, map[string]interface{}{"$user_id": "1"})
}
