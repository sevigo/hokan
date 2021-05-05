package watcher

import (
	notify "github.com/sevigo/notify/core"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

func (w *Watch) StartDirWatcher() {
	log.Debug("watcher.StartDirWatcher(): starting subscriber")
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirStart)
	wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Info("watcher.StartDirWatcher(): stream canceled")
			return
		case e := <-eventData:
			data, ok := e.Data.(*core.Directory)
			if !ok {
				log.Error("watcher.StartDirWatcher(): Can't add directory to the watch list")
				continue
			}
			w.notifier.StartWatching(data.Path, &notify.WatchingOptions{Rescan: true})
		}
	}
}

func (w *Watch) StartRescanWatcher() {
	log.Debug("watcher.StartRescanWatcher(): starting subscriber")
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirRescan)
	wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-eventData:
			w.notifier.RescanAll()
		}
	}
}

func (w *Watch) GetDirsToWatch() error {
	dirs, err := w.store.List(w.ctx)
	if err != nil {
		log.WithError(err).Print("Can't list all directories")
		return err
	}
	for _, dir := range dirs {
		log.Printf("watcher.GetDirsToWatch(): publish %#v", dir)
		err = w.event.Publish(w.ctx, &core.EventData{
			Type: core.WatchDirStart,
			Data: dir,
		})
		if err != nil {
			log.WithError(err).Print("Can't publish [WatchDirStart] event")
		}
	}
	return nil
}
