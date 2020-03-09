// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"github.com/sevigo/hokan/cmd/hokan/config"
)

// Injectors from wire.go:

func InitializeApplication(ctx context.Context, config2 config.Config) (application, error) {
	db, err := provideDatabase(config2)
	if err != nil {
		return application{}, err
	}
	directoryStore := provideDirectoryStore(db)
	eventCreator := provideEventCreator()
	logger := provideLogger(config2)
	server := apiServerProvider(directoryStore, eventCreator, logger)
	mainHealthzHandler := provideHealthz()
	handler := provideRouter(server, mainHealthzHandler)
	serverServer := provideServer(handler, config2)
	notifier := provideNotifier(ctx)
	watcher, err := provideWatcher(ctx, directoryStore, eventCreator, notifier)
	if err != nil {
		return application{}, err
	}
	fileStore := provideFileStore(db)
	target, err := provideTarget(ctx, fileStore, eventCreator)
	if err != nil {
		return application{}, err
	}
	mainApplication := newApplication(serverServer, directoryStore, watcher, target)
	return mainApplication, nil
}
