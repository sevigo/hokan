//+build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
)

func InitializeApplication(config config.Config) (application, error) {
	wire.Build(
		serverSet,
		storeSet,
		watcherSet,
		newApplication,
	)

	return application{}, nil
}
