package target

import (
	"context"
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
	fileStore      core.FileStore
	event          core.EventCreator
	register       map[string]core.TargetStorage
	registerStatus map[string]core.TargetStorageStatus
}

func New(ctx context.Context, fileStore core.FileStore, event core.EventCreator) (*Register, error) {
	log.Debug("target.New(): start")

	r := &Register{
		ctx:            ctx,
		fileStore:      fileStore,
		event:          event,
		register:       make(map[string]core.TargetStorage),
		registerStatus: make(map[string]core.TargetStorageStatus),
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
				r.rescanAllWatchedDirs()
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

func (r *Register) setTargetStatus(name string, status core.TargetStorageStatus) {
	r.Lock()
	defer r.Unlock()
	r.registerStatus[name] = status
}

func (r *Register) addTarget(name string, ts core.TargetStorage) {
	r.Lock()
	defer r.Unlock()
	r.register[name] = ts
	r.registerStatus[name] = core.TargetStorageOK
}

func (r *Register) getTarget(name string) core.TargetStorage {
	r.Lock()
	defer r.Unlock()
	return r.register[name]
}

func (r *Register) getTargetStatus(name string) core.TargetStorageStatus {
	r.Lock()
	defer r.Unlock()
	return r.registerStatus[name]
}
