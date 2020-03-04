package watcher

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/sevigo/hokan/pkg/core"
	notify "github.com/sevigo/notify/core"
)

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
			err := w.processWatchEvent(e)
			if err != nil {
				log.Err(err).Msg("Can't add/remove directory from the watch list")
			}
		}
	}
}

func (w *Watch) GetDirsToWatch() error {
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
				log.Err(err).Msg("Can't publish [WatchDirStart] event")
			}
		}
	}
	return nil
}

func (w *Watch) processWatchEvent(e *core.EventData) error {
	data, ok := e.Data.(*core.Directory)
	if !ok {
		return fmt.Errorf("can't convert EventData %#v", e)
	}

	switch e.Type {
	case core.WatchDirStart:
		w.notifier.StartWatching(data.Path, &notify.WatchingOptions{Rescan: true})
	case core.WatchDirStop:
		w.notifier.StopWatching(data.Path)
	}

	return nil
}
