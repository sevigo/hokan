package watcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/watcher/utils"
	"github.com/sevigo/notify/watcher"
)

func (w *Watch) StartFileWatcher() {
	ctx := w.ctx

	for {
		select {
		case <-ctx.Done():
			log.Printf("StartFileWatcher(): event stream canceled")
			return
		case ev := <-w.notifier.Event():
			log.WithFields(log.Fields{
				"event": watcher.ActionToString(ev.Action),
				"file":  ev.Path,
			}).Info("FileWatcher() event fired")
			// TODO: adapt ev.Action to core action
			err := w.publishFileChange(ev.Path)
			if err != nil {
				log.WithError(err).Error("watcher.StartFileWatcher(): Can't publish [FileAdded] event")
			}
		case err := <-w.notifier.Error():
			msg := fmt.Sprintf("[notifier] %q", err.Message)
			w.sse.PublishMessage(msg)
			log.WithField("level", err.Level).Error(msg)
		}
	}
}

func (w *Watch) publishFileChange(path string) error {
	checksum, info, err := utils.FileChecksumInfo(path)
	if err != nil {
		return err
	}
	// TODO maybe move this event to event.Publish ?
	w.sse.PublishMessage(fmt.Sprintf("[event] File %q added", path))
	return w.event.Publish(w.ctx, &core.EventData{
		Type: core.FileAdded,
		Data: core.File{
			Path:     path,
			Checksum: checksum,
			Info:     info,
		},
	})
}
