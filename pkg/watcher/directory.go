package watcher

import (
	"context"
	"sync"

	"github.com/prometheus/common/log"
	"github.com/sevigo/hokan/pkg/core"
)

var wg sync.WaitGroup

type Watch struct {
	ctx   context.Context
	event core.EventCreator
	store core.DirectoryStore
}

func New(ctx context.Context, dirStore core.DirectoryStore, event core.EventCreator) (*Watch, error) {
	w := &Watch{
		ctx:   ctx,
		event: event,
		store: dirStore,
	}
	wg.Add(1)
	go w.StartDirWatcher()
	wg.Wait()
	err := w.GetDirsToWatch()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Watch) StartDirWatcher() {
	log.Debugln("dir-watcher: starting subscriber")
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirStart)
	wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Debugln("dir-watcher: stream canceled")
			return
		case e := <-eventData:
			log.Debugf("dir-watcher: %#v", e.Data)
		}
	}
}

func (w *Watch) GetDirsToWatch() error {
	dirs, err := w.store.List(w.ctx)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir.Active {
			err = w.event.Publish(w.ctx, &core.EventData{
				Type: core.WatchDirStart,
				Data: dir,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
