package watcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/sevigo/hokan/pkg/core"
)

var wg sync.WaitGroup

type Watch struct {
	ctx      context.Context
	event    core.EventCreator
	store    core.DirectoryStore
	notifier core.Notifier
}

func New(ctx context.Context, dirStore core.DirectoryStore, event core.EventCreator, notifier core.Notifier) (*Watch, error) {
	log.Printf("watcher.New(): start")

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
	go w.StartFileWatcher()
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
			w.processWatchEvent(e)
		}
	}
}

func (w *Watch) GetDirsToWatch() error {
	fmt.Println(">>> watcher.GetDirsToWatch(): running publishers ...")
	dirs, err := w.store.List(w.ctx)
	if err != nil {
		log.Err(err).Msg("Can't list all directories")
		return err
	}
	for _, dir := range dirs {
		fmt.Printf(">>> GetDirsToWatch(): %#v\n", dir)
		if dir.Active {
			log.Printf("watcher.GetDirsToWatch(): publish %#v", dir)
			err = w.event.Publish(w.ctx, &core.EventData{
				Type: core.WatchDirStart,
				Data: dir,
			})
			if err != nil {
				log.Err(err).Msg("Can't publish [WatchDirStart] event")
			}
		}
	}
	return nil
}

func (w *Watch) processWatchEvent(e *core.EventData) error {
	data, ok := e.Data.(*core.Directory)
	if !ok {
		return fmt.Errorf("some error")
	}

	switch e.Type {
	case core.WatchDirStart:
		w.notifier.StartWatching(data.Path)
	case core.WatchDirStop:
		w.notifier.StopWatching(data.Path)
	}

	return nil
}
