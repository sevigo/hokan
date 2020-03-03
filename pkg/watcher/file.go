package watcher

import (
	"github.com/rs/zerolog/log"
	"github.com/sevigo/notify/watcher"
)

func (w *Watch) StartFileWatcher() {
	ctx := w.ctx

	for {
		select {
		case <-ctx.Done():
			log.Printf("dir-watcher: stream canceled")
			return
		case ev := <-w.notifier.Event():
			log.Printf("[EVENT] %s: %q", watcher.ActionToString(ev.Action), ev.Path)
		case err := <-w.notifier.Error():
			if err.Level == "ERROR" {
				log.Printf("[%s] %s", err.Level, err.Message)
			}
		}
	}
}
