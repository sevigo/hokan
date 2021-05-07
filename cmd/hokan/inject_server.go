package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/event"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/handler/events"
	"github.com/sevigo/hokan/pkg/handler/health"
	"github.com/sevigo/hokan/pkg/handler/web"
	"github.com/sevigo/hokan/pkg/server"
	"github.com/sevigo/hokan/pkg/sse"
)

type healthzHandler http.Handler

// wire set for loading the server.
var serverSet = wire.NewSet(
	apiServerProvider,
	eventsServerProvider,
	provideHealthz,
	provideWebHandler,
	provideRouter,
	provideServer,
	provideEventCreator,
	provideServerSideEventCreator,
)

func provideWebHandler() *web.Server {
	return web.New()
}

func apiServerProvider(
	fileStore core.FileStore,
	dirStore core.DirectoryStore,
	events core.EventCreator) *api.Server {
	return api.New(fileStore, dirStore, events)
}

func eventsServerProvider(sseCreator core.ServerSideEventCreator) *events.Server {
	return events.New(sseCreator)
}

func provideRouter(apiHandler *api.Server,
	webHandler *web.Server,
	serverEventsHandler *events.Server,
	healthz healthzHandler) http.Handler {
	r := chi.NewRouter()
	r.Mount("/healthz", healthz)
	r.Mount("/api", apiHandler.Handler())
	r.Mount("/sse", serverEventsHandler.Handler())
	r.Mount("/", webHandler.Handler())
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

func provideServerSideEventCreator(ctx context.Context, results chan core.BackupResult) core.ServerSideEventCreator {
	return sse.New(ctx, results)
}
