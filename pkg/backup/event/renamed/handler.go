package renamed

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

type fileRenamed struct {
	*core.EventHandler
}

func New(handler *core.EventHandler) core.EventProcessor {
	return &fileRenamed{handler}
}

func (f *fileRenamed) Process(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %v", e)
	}
	log.Printf("renamed.Process(): file %q is renamed to %q", file.OldPath, file.Path)
	f.Results <- core.BackupResult{
		Success: true,
		Message: core.BackupFilerenamedMessage,
	}
	return nil
}
