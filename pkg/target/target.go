package target

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/void"
)

// All known targets
var targets = map[string]core.TargetFactory{
	void.TargetName: void.New,
}

type register struct {
	ctx context.Context
	sync.Mutex
	fileStore core.FileStore
	event     core.EventCreator
	register  map[string]core.TargetStorage
}

func New(ctx context.Context, fileStore core.FileStore, event core.EventCreator) (*register, error) {
	r := &register{
		ctx:       ctx,
		fileStore: fileStore,
		event:     event,
		register:  make(map[string]core.TargetStorage),
	}
	r.InitTargets()
	go r.StartFileAddedWatcher()
	return r, nil
}

func (r *register) InitTargets() {
	for name, target := range targets {
		ts, err := target(r.fileStore)
		if err != nil {
			log.Err(err).Msg("Can't create new target storage")
			continue
		}
		r.addTarget(name, ts)
	}
}

func (r *register) addTarget(name string, ts core.TargetStorage) {
	r.Lock()
	defer r.Unlock()
	r.register[name] = ts
}

func (r *register) StartFileAddedWatcher() {
	log.Printf("target.StartFileChangeWatcher(): starting subscriber")
	ctx := r.ctx
	eventData := r.event.Subscribe(ctx, core.FileAdded)

	for {
		select {
		case <-ctx.Done():
			log.Printf("file-watcher: stream canceled")
			return
		case e := <-eventData:
			log.Printf("file-watcher: %#v", e.Data)
			// err := w.processWatchEvent(e)
			// if err != nil {
			// 	log.Err(err).Msg("Can't add/remove directory from the watch list")
			// }
		}
	}
}
