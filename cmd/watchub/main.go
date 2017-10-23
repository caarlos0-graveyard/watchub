package main

import (
	"net/http"
	"time"

	"github.com/caarlos0/watchub/oauth"
	"github.com/caarlos0/watchub/postgres"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/caarlos0/watchub/config"
)

func main() {
	log.SetHandler(logfmt.Default)
	log.SetLevel(log.InfoLevel)
	log.Info("starting up...")

	var config = config.Get()
	var db = postgres.Connect(config.DatabaseURL)
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Fatal("failed to close database connections")
		}
	}()

	var session = sessions.NewCookieStore([]byte(config.SessionSecret))
	var oauth = oauth.New(config)

	var mux = mux.NewRouter()
	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// prometheus stuff
	mux.Handle("/metrics", promhttp.Handler())

	var handler = context.ClearHandler(httplog.New(mux))

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
