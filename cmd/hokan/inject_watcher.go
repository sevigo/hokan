package main

import (
	"context"

	"github.com/google/wire"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/watcher"

	"github.com/sevigo/notify"
	notifywatcher "github.com/sevigo/notify/watcher"
)

var watcherSet = wire.NewSet(
	provideWatcher,
	provideNotifier,
)

func provideWatcher(ctx context.Context,
	dirStore core.DirectoryStore,
	event core.EventCreator,
	notifier core.Notifier,
	sse core.ServerSideEventCreator) (core.Watcher, error) {
	return watcher.New(ctx, dirStore, event, notifier, sse)
}

func provideNotifier(ctx context.Context) core.Notifier {
	return notify.Setup(ctx, &notifywatcher.Options{})
}
