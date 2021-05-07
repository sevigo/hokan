package watcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/watcher/utils"
	"github.com/sevigo/notify/event"
	"github.com/sevigo/notify/watcher"
)

func (w *Watch) StartFileWatcher() {
	log.Printf("watcher.StartFileWatcher(): starting")
	ctx := w.ctx

	for {
		select {
		case <-ctx.Done():
			log.Printf("watcher.StartFileWatcher(): event-stream canceled")
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
		case e := <-w.notifier.Error():
			w.printNotifyMessage(e)
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

// TODO: find a better place for me
func (w *Watch) printNotifyMessage(e event.Error) {
	msg := fmt.Sprintf("[notifier] %q", e.Message)
	switch e.Level {
	case "DEBUG":
		log.WithField("level", e.Level).Debug(msg)
	case "INFO":
		log.WithField("level", e.Level).Info(msg)
	case "WARN", "WARNING":
		log.WithField("level", e.Level).Warn(msg)
	case "ERROR":
		log.WithField("level", e.Level).Error(msg)
	case "CRITICAL":
		log.WithField("level", e.Level).Fatal(msg)
	default:
		log.Print(msg)
	}
	w.sse.PublishMessage(msg)
}
