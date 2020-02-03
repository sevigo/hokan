package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/handler/health"
	"github.com/sevigo/hokan/pkg/server"
)

type healthzHandler http.Handler

// wire set for loading the server.
var serverSet = wire.NewSet(
	provideHealthz,
	provideRouter,
	provideServer,
)

func provideRouter(healthz healthzHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/healthz", healthz)
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
