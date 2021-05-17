package removed

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

type fileRemoved struct {
	*core.EventHandler
}

func New(handler *core.EventHandler) core.EventProcessor {
	return &fileRemoved{handler}
}

func (f *fileRemoved) Name() string {
	return core.EventToString(core.FileRemoved)
}

func (f *fileRemoved) Process(e *core.EventData) error {
	fmt.Printf(">>> Process(): %+v\n", e)
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %+v", e)
	}
	backupName := f.Backup.Name()
	log.Printf("removed.Process(): file %q is deleted from backup [%s]", file.Path, backupName)
	err := f.FileStore.Delete(f.Ctx, backupName, &file)
	if err == nil {
		f.Results <- core.BackupResult{
			Success: true,
			Message: core.BackupFileDeletedMessage,
		}
	} else {
		f.Results <- core.BackupResult{
			Success: false,
			Error:   err,
		}
	}
	log.Printf("removed.Process(): Delete: %v", err)
	return err
}
