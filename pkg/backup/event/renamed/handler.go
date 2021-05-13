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
	log.Printf("event.Process(): file rename %+v\n", file)
	return nil
}
