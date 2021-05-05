//+build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
)

// for more info see: https://github.com/google/wire
func InitializeApplication(ctx context.Context, config config.Config) (application, error) {
	wire.Build(
		serverSet,
		storeSet,
		watcherSet,
		backupSet,
		newApplication,
	)

	return application{}, nil
}
