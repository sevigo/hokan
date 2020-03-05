package watcher

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return "", "", err
	}

	if info.IsDir() {
		return "", "", fmt.Errorf("not a file")
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", "", err
	}

	infoJson, err := json.Marshal(FileInfo{info})
	if _, err := io.Copy(h, f); err != nil {
		return "", "", err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))

	return sum, string(infoJson), nil
}
