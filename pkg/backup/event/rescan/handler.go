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

func (f *dirRescan) Name() string {
	return core.EventToString(core.WatchDirRescan)
}

func (d *dirRescan) Process(event *core.EventData) error {
	log.Infof("rescan.Process() for dir: %+v", event.Data)
	return nil
}
