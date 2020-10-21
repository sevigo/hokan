package watcher

import (
	"context"
	"sync"

	"github.com/sevigo/hokan/pkg/core"
)

var wg sync.WaitGroup

type Watch struct {
	ctx      context.Context
	event    core.EventCreator
	store    core.DirectoryStore
	notifier core.Notifier
	catalog  []*core.Directory
	sse      core.ServerSideEventCreator
}

func New(ctx context.Context,
	dirStore core.DirectoryStore,
	event core.EventCreator,
	notifier core.Notifier,
	sse core.ServerSideEventCreator) (*Watch, error) {
	w := &Watch{
		ctx:      ctx,
		event:    event,
		store:    dirStore,
		notifier: notifier,
		sse:      sse,
	}
	wg.Add(2) //nolint:gomnd
	go w.StartDirWatcher()
	go w.StartRescanWatcher()
	wg.Wait()
	err := w.GetDirsToWatch()
	if err != nil {
		return nil, err
	}

	dirs, err := dirStore.List(ctx)
	if err != nil {
		return nil, err
	}
	w.catalog = dirs

	go w.StartFileWatcher()
	return w, nil
}
