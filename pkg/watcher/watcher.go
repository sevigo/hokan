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
}

func New(ctx context.Context, dirStore core.DirectoryStore, event core.EventCreator, notifier core.Notifier) (*Watch, error) {
	w := &Watch{
		ctx:      ctx,
		event:    event,
		store:    dirStore,
		notifier: notifier,
	}
	wg.Add(1) //nolint:gomnd
	go w.StartDirWatcher()
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
