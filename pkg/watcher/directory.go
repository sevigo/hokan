package watcher

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/sevigo/hokan/pkg/core"
)

var wg sync.WaitGroup

type Watch struct {
	ctx   context.Context
	event core.EventCreator
	store core.DirectoryStore
}

func New(ctx context.Context, dirStore core.DirectoryStore, event core.EventCreator) (*Watch, error) {
	log.Printf("watcher.New(): start")

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
	log.Printf("watcher.StartDirWatcher(): starting subscriber")
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirStart)
	wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("dir-watcher: stream canceled")
			return
		case e := <-eventData:
			log.Printf("dir-watcher: %#v", e.Data)
		}
	}
}

func (w *Watch) GetDirsToWatch() error {
	log.Printf("watcher.GetDirsToWatch(): running publishers")
	dirs, err := w.store.List(w.ctx)
	if err != nil {
		log.Err(err).Msg("Can't list all directories")
		return err
	}
	for _, dir := range dirs {
		if dir.Active {
			log.Printf("watcher.GetDirsToWatch(): publish %#v", dir)
			err = w.event.Publish(w.ctx, &core.EventData{
				Type: core.WatchDirStart,
				Data: dir,
			})
			if err != nil {
				// return err
				log.Err(err).Msg("Can't publish [WatchDirStart] event")
			}
		}
	}
	return nil
}
