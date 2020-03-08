package watcher

import (
	"fmt"

	"github.com/rs/zerolog/log"
	notify "github.com/sevigo/notify/core"

	"github.com/sevigo/hokan/pkg/core"
)

func (w *Watch) StartDirWatcher() {
	log.Debug().Msg("watcher.StartDirWatcher(): starting subscriber")
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirStart)
	wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("dir-watcher: stream canceled")
			return
		case e := <-eventData:
			data, ok := e.Data.(*core.Directory)
			if !ok {
				log.Error().Msg("Can't add directory to the watch list")
			}
			w.notifier.StartWatching(data.Path, &notify.WatchingOptions{Rescan: true})
		}
	}
}

func (w *Watch) StartRescanWatcher() {
	log.Debug().Msg("watcher.StartRescanWatcher(): starting subscriber")
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirRescan)
	wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-eventData:
			fmt.Println(">>> received event WatchDirRescan")
			w.notifier.RescanAll()
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
