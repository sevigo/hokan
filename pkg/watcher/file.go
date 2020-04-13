package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
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
			log.WithField("level", err.Level).Errorf("[notifier] %q", err.Message)
		}
	}
}

func (w *Watch) publishFileChange(path string) error {
	var targets []string
	for _, dir := range w.catalog {
		if strings.Contains(path, dir.Path) {
			targets = dir.Targets
		}
	}

	checksum, info, err := fileChecksumInfo(path)
	if err != nil {
		return err
	}
	return w.event.Publish(w.ctx, &core.EventData{
		Type: core.FileAdded,
		Data: core.File{
			Path:     path,
			Checksum: checksum,
			Info:     info,
			Targets:  targets,
		},
	})
}

func fileChecksumInfo(path string) (string, *core.FileInfo, error) {
	f, erro := os.Open(path)
	if erro != nil {
		return "", nil, erro
	}
	defer f.Close()
	info, errs := f.Stat()
	if errs != nil {
		return "", nil, errs
	}

	if info.IsDir() {
		return "", nil, fmt.Errorf("not a file")
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil, err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, &core.FileInfo{info}, nil
}
