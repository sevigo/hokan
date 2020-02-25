//+build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
)

func InitializeApplication(ctx context.Context, config config.Config) (application, error) {
	wire.Build(
		serverSet,
		storeSet,
		watcherSet,
		newApplication,
	)

	return application{}, nil
}
