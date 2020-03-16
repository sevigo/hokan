package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/event"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/handler/health"
	"github.com/sevigo/hokan/pkg/server"
)

type healthzHandler http.Handler

// wire set for loading the server.
var serverSet = wire.NewSet(
	apiServerProvider,
	provideHealthz,
	provideRouter,
	provideServer,
	provideEventCreator,
	provideLogger,
)

func apiServerProvider(dirStore core.DirectoryStore, events core.EventCreator, targets core.TargetRegister, logger *logrus.Logger) *api.Server {
	return api.New(dirStore, events, targets, logger)
}

func provideLogger(config config.Config) *logrus.Logger {
	l := logrus.StandardLogger()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: config.Logging.Color,
	})
	l.SetReportCaller(false)
	if config.Logging.Debug {
		l.SetLevel(logrus.DebugLevel)
	}
	l.SetOutput(os.Stdout)
	return l
}

func provideRouter(apiHandler *api.Server, healthz healthzHandler) http.Handler {
	r := chi.NewRouter()
	r.Mount("/healthz", healthz)
	r.Mount("/api", apiHandler.Handler())
	return r
}

func provideHealthz() healthzHandler {
	v := health.New()
	return healthzHandler(v)
}

func provideServer(handler http.Handler, config config.Config) *server.Server {
	return &server.Server{
		Addr:    config.Server.Addr,
		Host:    config.Server.Host,
		Handler: handler,
	}
}

func provideEventCreator() core.EventCreator {
	return event.New(event.Config{})
}
