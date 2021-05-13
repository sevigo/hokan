package rescan

import (
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

type dirRescan struct {
	*core.EventHandler
}

func New(handler *core.EventHandler) core.EventProcessor {
	return &dirRescan{handler}
}

func (d *dirRescan) Process(event *core.EventData) error {
	log.Infof("Process dir rescan: %#v", event.Data)
	return nil
}
