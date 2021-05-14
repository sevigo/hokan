package event

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/backup/event/added"
	"github.com/sevigo/hokan/pkg/backup/event/changed"
	"github.com/sevigo/hokan/pkg/backup/event/removed"
	"github.com/sevigo/hokan/pkg/backup/event/renamed"
	"github.com/sevigo/hokan/pkg/backup/event/rescan"
	"github.com/sevigo/hokan/pkg/core"
)

var eventFactory = map[core.EventType]core.EventProcessrFactory{
	core.FileAdded:      added.New,
	core.FileChanged:    changed.New,
	core.FileRemoved:    removed.New,
	core.WatchDirRescan: rescan.New,
	core.FileRenamed:    renamed.New,
}

func InitHanler(ctx context.Context, events core.EventCreator, backup core.Backup, fileStore core.FileStore) {
	handler := &core.EventHandler{
		Ctx:       ctx,
		Events:    events,
		Backup:    backup,
		FileStore: fileStore,
	}

	for eventType, factory := range eventFactory {
		processor := factory(handler)
		event := events.Subscribe(ctx, eventType)
		go listener(ctx, event, processor)
	}
}

func listener(ctx context.Context, event <-chan *core.EventData, processor core.EventProcessor) {
	for {
		select {
		case <-ctx.Done():
			log.Info("event.Listener(): event stream canceled")
			return
		case e := <-event:
			log.Infof("backup.Listener(): event [%s] is triggerd", core.EventToString(e.Type))
			err := processor.Process(e)
			if err != nil {
				log.WithError(err).
					WithField("event", core.EventToString(e.Type)).
					Error("event.Listener(): can't process event")
			}
		}
	}
}
