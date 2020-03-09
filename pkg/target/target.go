package target

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/minio"
	"github.com/sevigo/hokan/pkg/target/void"
)

// All known targets
var targets = map[string]core.TargetFactory{
	void.TargetName:  void.New,
	minio.TargetName: minio.New,
}

type Register struct {
	ctx context.Context
	sync.Mutex
	fileStore core.FileStore
	event     core.EventCreator
	register  map[string]core.TargetStorage
}

func New(ctx context.Context, fileStore core.FileStore, event core.EventCreator) (*Register, error) {
	log.Debug("target.New(): start")

	r := &Register{
		ctx:       ctx,
		fileStore: fileStore,
		event:     event,
		register:  make(map[string]core.TargetStorage),
	}
	r.InitTargets(ctx)
	go r.StartFileAddedWatcher()
	return r, nil
}

func (r *Register) InitTargets(ctx context.Context) {
	for name, target := range targets {
		go r.initWithRetry(ctx, name, target)
	}
}

func (r *Register) initWithRetry(ctx context.Context, name string, target core.TargetFactory) {
	// first run
	ts, err := target(ctx, r.fileStore)
	if err == nil {
		r.addTarget(name, ts)
		return
	}

	// retry
	counter := 1
	mod := 10
	ticker := time.NewTicker(time.Duration(mod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("stream canceled")
			return
		case <-ticker.C:
			log.WithFields(log.Fields{
				"target":         name,
				"retry-counter":  counter,
				"offset-seconds": mod,
			}).Error("retry target setup")
			ts, err := target(ctx, r.fileStore)
			if err == nil {
				r.addTarget(name, ts)
				errEvent := r.event.Publish(ctx, &core.EventData{Type: core.WatchDirRescan})
				if errEvent != nil {
					log.WithError(errEvent).Error("Can't publish [FileAdded] event")
				}
				return
			}
			log.WithError(err).Error("Can't create new target storage")
			if counter%10 == 0 {
				mod *= 2
				ticker = time.NewTicker(time.Duration(mod) * time.Second)
			}
			counter++
		}
	}
}

func (r *Register) StartFileAddedWatcher() {
	ctx := r.ctx
	eventData := r.event.Subscribe(ctx, core.FileAdded)

	for {
		select {
		case <-ctx.Done():
			log.Printf("file-watcher: stream canceled")
			return
		case e := <-eventData:
			// log.Debugf("StartFileAddedWatcher(): %#v\n", e.Data)
			err := r.processFileAddedEvent(e)
			if err != nil {
				log.WithError(err).Error("can't send the file to the target storage")
			}
		}
	}
}

func (r *Register) processFileAddedEvent(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %v", e)
	}
	for _, target := range file.Targets {
		if ts := r.getTarget(target); ts != nil {
			err := ts.Save(r.ctx, &file)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"target": target,
					"file":   file.Path,
				}).Error("can't save the file to the target storage")
			}
		}
	}
	return nil
}

func (r *Register) addTarget(name string, ts core.TargetStorage) {
	r.Lock()
	defer r.Unlock()
	r.register[name] = ts
}

func (r *Register) getTarget(name string) core.TargetStorage {
	r.Lock()
	defer r.Unlock()
	return r.register[name]
}
