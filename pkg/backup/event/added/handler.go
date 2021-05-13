package added

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/backup/utils"
	"github.com/sevigo/hokan/pkg/core"
)

type fileAdded struct {
	*core.EventHandler
}

func New(handler *core.EventHandler) core.EventProcessor {
	return &fileAdded{handler}
}

func (f *fileAdded) Process(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %v", e)
	}
	log.Printf("watcher.processFileAddedEvent(): event [%s] was fired\n", core.EventToString(e.Type))
	f.saveFileToBackup(&file)
	return nil
}

func (f *fileAdded) saveFileToBackup(file *core.File) {
	storedFile, err := f.FileStore.Find(f.Ctx, f.Backup.Name(), &core.FileSearchOptions{
		FilePath: file.Path,
	})
	if errors.Is(err, core.ErrFileNotFound) || utils.FileHasChanged(file, storedFile) {
		go f.Backup.Save(f.Ctx, f.Results, file, &core.BackupOperationOptions{})
		return
	}
	log.Printf("event.saveFileToBackup(): ignore %q %s\n", file.Path, core.BackupNoChangeMessage)
	f.Results <- core.BackupResult{
		Success: true,
		Message: core.BackupNoChangeMessage,
	}
}
