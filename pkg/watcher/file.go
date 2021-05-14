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
			switch ev.Action {
			case event.FileAdded: //, event.FileModified:
				err := w.publishFileAdded(ev.Path)
				if err != nil {
					log.WithError(err).Error("watcher.StartFileWatcher(): Can't publish [FileAdded] event")
				}
			case event.FileModified:
				err := w.publishFileChanged(ev.Path)
				if err != nil {
					log.WithError(err).Error("watcher.StartFileWatcher(): Can't publish [FileChanged] event")
				}
			case event.FileRenamedNewName:
				err := w.publishFileRenamed(ev.Path, ev.AdditionalInfo.OldName)
				if err != nil {
					log.WithError(err).Error("watcher.StartFileWatcher(): Can't publish [FileRenamed] event")
				}
			case event.FileRemoved:
				w.event.Publish(w.ctx, &core.EventData{
					Type: core.FileRemoved,
					Data: core.File{
						Path: ev.Path,
					},
				})

			default:
				log.Infof("ignoring this event")
			}
		case e := <-w.notifier.Error():
			w.printNotifyMessage(e)
		}
	}
}

func (w *Watch) publishFileRenamed(newPath, oldPath string) error {
	if oldPath == "" {
		return fmt.Errorf("publishFileRenamed(): old path can't be empty")
	}
	checksum, info, err := utils.FileChecksumInfo(newPath)
	if err != nil {
		return err
	}
	// TODO maybe move this event to event.Publish ?
	w.sse.PublishMessage(fmt.Sprintf("[event] File %q was renamed to %q", oldPath, newPath))
	return w.event.Publish(w.ctx, &core.EventData{
		Type: core.FileRenamed,
		Data: core.File{
			Path:     newPath,
			OldPath:  oldPath,
			Checksum: checksum,
			Info:     info,
		},
	})
}

func (w *Watch) publishFileAdded(path string) error {
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

func (w *Watch) publishFileChanged(path string) error {
	checksum, info, err := utils.FileChecksumInfo(path)
	if err != nil {
		return err
	}
	w.sse.PublishMessage(fmt.Sprintf("[event] File %q changed", path))
	return w.event.Publish(w.ctx, &core.EventData{
		Type: core.FileChanged,
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
