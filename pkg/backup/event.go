package backup

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/backup/utils"
)

type Watcher struct {
	ctx       context.Context
	events    core.EventCreator
	backup    core.Backup
	fileStore core.FileStore
	// write results of any operation here, this will be propagated to the UI
	Results chan core.BackupResult
}

func (w *Watcher) FileAdded() {
	log.Printf("watcher.FileAdded(): starting")
	ctx := w.ctx
	eventData := w.events.Subscribe(ctx, core.FileAdded)

	for {
		select {
		case <-ctx.Done():
			log.Info("backup.FileAdded(): event stream canceled")
			return
		case e := <-eventData:
			err := w.processFileAddedEvent(e)
			if err != nil {
				log.WithError(err).Error("backup.FileAdded(): can't send the file to the target storage")
			}
		}
	}
}

func (w *Watcher) processFileAddedEvent(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %v", e)
	}
	log.Printf("watcher.processFileAddedEvent(): event [%s] was fired\n", core.EventToString(e.Type))
	w.saveFileToBackup(&file)
	return nil
}

func (w *Watcher) saveFileToBackup(file *core.File) {
	storedFile, err := w.fileStore.Find(w.ctx, &core.FileSearchOptions{
		FilePath: file.Path,
	})
	if errors.Is(err, core.ErrFileNotFound) || utils.FileHasChanged(file, storedFile) {
		go w.backup.Save(w.ctx, w.Results, file, &core.BackupOperationOptions{})
		return
	}
	log.Printf("watcher.saveFileToBackup(): ignore %q %s\n", file.Path, core.BackupNoChangeMessage)
	w.Results <- core.BackupResult{
		Success: true,
		Message: core.BackupNoChangeMessage,
	}
}

func (w *Watcher) rescanAllWatchedDirs() {
	log.Info("backup.rescanAllWatchedDirs(): publish [WatchDirRescan] event")
	err := w.events.Publish(w.ctx, &core.EventData{Type: core.WatchDirRescan})
	if err != nil {
		log.WithError(err).Error("Can't publish [WatchDirRescan] event")
	}
}
