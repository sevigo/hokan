//+build wireinject
package main

import (
	"net/http"

	"github.com/google/wire"

	"github.com/seviko/hokan/cmd/hokan/config"
	"github.com/seviko/hokan/pkg/handler/health"
)

type healthzHandler http.Handler

// wire set for loading the server.
var serverSet = wire.NewSet(
	provideHealthz,
)

func provideHealthz() healthzHandler {
	v := health.New()
	return healthzHandler(v)
}

func InitializeApplication(config config.Config) (application, error) {
	wire.Build(serverSet)

	return application{}, nil
}
