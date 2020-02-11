package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/wire"

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
	api.New,
	provideHealthz,
	provideRouter,
	provideServer,
	provideEventCreator,
)

func provideRouter(api api.Server, healthz healthzHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/healthz", healthz)
	r.Mount("/api", api.Handler())
	return r
}

func provideHealthz() healthzHandler {
	v := health.New()
	return healthzHandler(v)
}

func provideServer(handler *chi.Mux, config config.Config) *server.Server {
	return &server.Server{
		Addr:    config.Server.Addr,
		Host:    config.Server.Host,
		Handler: handler,
	}
}

func provideEventCreator(config event.Config) core.EventCreator {
	return event.New(event.Config{})
}
