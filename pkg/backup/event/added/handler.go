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

func (f *fileAdded) Name() string {
	return core.EventToString(core.FileAdded)
}

func (f *fileAdded) Process(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %v", e)
	}
	log.Printf("added.Process(): for %q was fired ", file.Path)
	f.saveFileToBackup(&file)
	return nil
}

func (f *fileAdded) saveFileToBackup(file *core.File) {
	storedFile, err := f.FileStore.Find(f.Ctx, f.Backup.Name(), &core.FileSearchOptions{
		FilePath: file.Path,
	})
	// file is new and is not found in the backup
	if errors.Is(err, core.ErrFileNotFound) {
		go f.Backup.Save(f.Ctx, f.Results, file, &core.BackupOperationOptions{})
		return
	}

	// something else is wrong
	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{
				"backup": f.Backup.Name(),
				"file":   file.Path,
			}).
			Error("Can't find file in the local storage")
		f.Results <- core.BackupResult{
			Success: false,
			Message: "Can't find file in the local storage",
			Error:   err,
		}
	}

	if utils.FileHasChanged(file, storedFile) {
		log.Printf("added.saveFileToBackup(): file %q has changed since the last update", file.Path)
		go f.Backup.Save(f.Ctx, f.Results, file, &core.BackupOperationOptions{})
		return
	}

	log.Printf("added.saveFileToBackup(): file %q hasn't changed since the last update", file.Path)
	f.Results <- core.BackupResult{
		Success: true,
		Message: core.BackupNoChangeMessage,
	}
}
