package main

import (
	"context"

	"github.com/google/wire"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/watcher"
)

var watcherSet = wire.NewSet(
	provideWatcher,
)

func provideWatcher(ctx context.Context, dirStore core.DirectoryStore, event core.EventCreator) (core.Watcher, error) {
	return watcher.New(ctx, dirStore, event)
}
