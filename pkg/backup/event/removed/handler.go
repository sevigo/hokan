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

func (f *fileRemoved) Process(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %+v", e)
	}
	log.Printf("removed.Process(): file %q is renamed", file.Path)
	f.Results <- core.BackupResult{
		Success: true,
		Message: core.BackupFileDeletedMessage,
	}
	return nil
}
