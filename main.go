package main

//_ "net/http/pprof"

import (
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/controllers"
	"github.com/caarlos0/watchub/datastore/database"
	"github.com/caarlos0/watchub/oauth"
	"github.com/caarlos0/watchub/scheduler"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: refactor this.
// nolint: funlen
func main() {
	log.SetHandler(logfmt.Default)
	log.SetLevel(log.InfoLevel)
	log.Info("starting up...")

	var config = config.Get()
	var db = database.Connect()
	defer func() { _ = db.Close() }()
	var store = database.NewDatastore(db)

	// oauth
	var session = sessions.NewCookieStore([]byte(config.SessionSecret))
	session.Options = &sessions.Options{
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	var oauth = oauth.New(config)
	var loginCtrl = controllers.NewLogin(config, session, oauth, store)

	// schedulers
	var sch = scheduler.New(config, store, oauth, session)
	sch.Start()
	defer sch.Stop()

	// routes
	var mux = mux.NewRouter()
	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)
	mux.Methods(http.MethodGet).Path("/").HandlerFunc(
		controllers.NewIndex(config, session, store).Handler,
	)
	mux.Methods(http.MethodGet).Path("/donate").HandlerFunc(
		controllers.NewDonate(config, session).Handler,
	)
	mux.Methods(http.MethodGet).Path("/contact").HandlerFunc(
		controllers.NewContact(config, session).Handler,
	)
	mux.Methods(http.MethodGet).Path("/schedule").HandlerFunc(
		controllers.NewSchedule(config, session, store).Handler,
	)
	mux.Methods(http.MethodGet).Path("/login").HandlerFunc(
		loginCtrl.Handler,
	)
	mux.Methods(http.MethodGet).Path("/login/callback").HandlerFunc(
		loginCtrl.CallbackHandler,
	)
	mux.Path("/logout").HandlerFunc(
		controllers.NewLogout(config, session).Handler,
	)

	// prometheus stuff
	var requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "watchub",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "total requests",
	}, []string{"code", "method"})
	var responseObserver = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "watchub",
		Subsystem: "http",
		Name:      "responses",
		Help:      "response times and counts",
	}, []string{"code", "method"})

	prometheus.MustRegister(scheduler.TimeGauge)
	prometheus.MustRegister(scheduler.ErrorGauge)
	mux.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	mux.PathPrefix("/debug").Handler(http.DefaultServeMux)

	var handler = context.ClearHandler(
		promhttp.InstrumentHandlerDuration(
			responseObserver,
			promhttp.InstrumentHandlerCounter(
				requestCounter,
				mux,
			),
		),
	)

	var server = &http.Server{
		Handler:      handler,
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithField("addr", server.Addr).Info("started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start up server")
	}
}
