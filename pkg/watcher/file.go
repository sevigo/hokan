package watcher

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sevigo/notify/watcher"

	"github.com/sevigo/hokan/pkg/core"
)

func (w *Watch) StartFileWatcher() {
	ctx := w.ctx

	for {
		select {
		case <-ctx.Done():
			log.Printf("file-watcher: stream canceled")
			return
		case ev := <-w.notifier.Event():
			log.Printf("[EVENT] %s: %q", watcher.ActionToString(ev.Action), ev.Path)
			// TODO: adapt ev.Action to core action
			err := w.publishFileChange(ev.Path)
			if err != nil {
				log.Err(err).Msg("Can't publish [FileAdded] event")
			}
		case err := <-w.notifier.Error():
			if err.Level == "ERROR" {
				log.Printf("[%s] %s", err.Level, err.Message)
			}
		}
	}
}

func (w *Watch) publishFileChange(path string) error {
	var targets []string
	for _, dir := range w.catalog {
		if strings.Contains(path, dir.Path) {
			targets = dir.Target
		}
	}

	return w.event.Publish(w.ctx, &core.EventData{
		Type: core.FileAdded,
		Data: core.File{
			Path:    path,
			Targets: targets,
		},
	})
}
