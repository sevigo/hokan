package target

import (
	"context"
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/local"
	"github.com/sevigo/hokan/pkg/target/minio"
	"github.com/sevigo/hokan/pkg/target/void"
)

// All known targets
var targets = map[string]core.TargetFactory{
	local.TargetName: local.New,
	minio.TargetName: minio.New,
	void.TargetName:  void.New,
}

type Register struct {
	ctx context.Context
	sync.Mutex
	fileStore      core.FileStore
	configStore    core.ConfigStore
	event          core.EventCreator
	register       map[string]core.TargetStorage
	registerStatus map[string]core.TargetStorageStatus
}

func New(ctx context.Context, fileStore core.FileStore, configStore core.ConfigStore, event core.EventCreator) (core.TargetRegister, error) {
	log.Debug("target.New(): start")

	r := &Register{
		ctx:            ctx,
		fileStore:      fileStore,
		configStore:    configStore,
		event:          event,
		register:       make(map[string]core.TargetStorage),
		registerStatus: make(map[string]core.TargetStorageStatus),
	}
	r.initTargets(ctx)
	go r.StartFileAddedWatcher()
	return r, nil
}

func (r *Register) AllTargets() map[string]core.TargetFactory {
	return targets
}

func (r *Register) GetTarget(name string) core.TargetStorage {
	r.Lock()
	defer r.Unlock()
	return r.register[name]
}

func (r *Register) initTargets(ctx context.Context) {
	for name, target := range targets {
		go r.initWithRetry(ctx, name, target)
	}
}

func (r *Register) initWithRetry(ctx context.Context, name string, target core.TargetFactory) {
	// first run
	conf, err := r.GetConfig(ctx, name)
	if err != nil {
		log.WithFields(log.Fields{
			"target": name,
		}).WithError(err).Fatal("can't get configuration for the target storage")
		return
	}
	ts, err := target(ctx, r.fileStore, *conf)
	if err == nil {
		r.addTarget(name, ts)
		return
	}
	if errors.Is(err, core.ErrTargetNotActive) {
		log.WithFields(log.Fields{
			"target": name,
		}).WithError(err).Error("stop retry")
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
			ts, err := target(ctx, r.fileStore, *conf)
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

func (r *Register) getTargetStatus(name string) core.TargetStorageStatus {
	r.Lock()
	defer r.Unlock()
	return r.registerStatus[name]
}
