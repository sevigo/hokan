package watcher

import (
	"crypto/sha256"
	"encoding/json"
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

type FileInfo struct {
	os.FileInfo
}

func (f FileInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Name":    f.Name(),
		"Size":    f.Size(),
		"Mode":    f.Mode(),
		"ModTime": f.ModTime(),
	})
}

func fileChecksumInfo(path string) (string, string, error) {
	f, erro := os.Open(path)
	if erro != nil {
		return "", "", erro
	}
	defer f.Close()
	info, errs := f.Stat()
	if errs != nil {
		return "", "", errs
	}

	if info.IsDir() {
		return "", "", fmt.Errorf("not a file")
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", "", err
	}

	infoJSON, errj := json.Marshal(FileInfo{info})
	if errj != nil {
		return "", "", errj
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, string(infoJSON), nil
}
