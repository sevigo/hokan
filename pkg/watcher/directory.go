package watcher

import (
	"context"

	"github.com/prometheus/common/log"
	"github.com/sevigo/hokan/pkg/core"
)

type Watch struct {
	ctx   context.Context
	event core.EventCreator
}

func New(ctx context.Context, event core.EventCreator) *Watch {
	w := &Watch{
		ctx:   ctx,
		event: event,
	}
	go w.Start()
	return w
}

func (w *Watch) Start() {
	ctx := w.ctx
	eventData := w.event.Subscribe(ctx, core.WatchDirStart)

	for {
		select {
		case <-ctx.Done():
			log.Debugln("watcher: stream canceled")
			return
		case e := <-eventData:
			log.Debugf("watcher: %#v", e.Data)
		}
	}
}
