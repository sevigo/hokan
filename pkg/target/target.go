package target

import (
	"context"
	"errors"
	"fmt"
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

func New(
	ctx context.Context,
	fileStore core.FileStore,
	configStore core.ConfigStore,
	event core.EventCreator) (core.TargetRegister, error) {
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

func (r *Register) AllTargets() map[string]core.Target {
	result := map[string]core.Target{}
	for name, status := range r.registerStatus {
		t := r.GetTarget(name)
		if t == nil {
			continue
		}
		result[name] = core.Target{
			Name:   name,
			Status: status,
			Info:   t.Info(r.ctx),
		}
	}
	return result
}

func (r *Register) GetTarget(name string) core.TargetStorage {
	r.Lock()
	defer r.Unlock()
	return r.register[name]
}

func (r *Register) initTargets(ctx context.Context) {
	for name := range targets {
		err := r.initTarget(ctx, name)
		if errors.Is(err, core.ErrTargetNotActive) {
			log.WithError(err).
				WithField("target", name).
				Error("initTargets(): ignore the target for now")
			continue
		}
		if err != nil {
			log.WithError(err).Error("initTargets()")
			go r.initWithRetry(ctx, name)
		}
	}
}

func (r *Register) initTarget(ctx context.Context, name string) error {
	logger := log.WithFields(log.Fields{
		"target": name,
	})

	target, ok := targets[name]
	if !ok {
		return fmt.Errorf("initTarget(): target %q not found", name)
	}

	conf, err := r.GetConfig(ctx, name)
	if err != nil {
		logger.WithError(err).Fatal("can't get configuration for the target storage")
		return err
	}
	ts, err := target(ctx, r.fileStore, *conf)
	if err != nil {
		return err
	}

	r.addTarget(name, ts)
	return nil
}

func (r *Register) initWithRetry(ctx context.Context, name string) {
	counter := 1
	mod := 10
	ticker := time.NewTicker(time.Duration(mod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("initWithRetry(): stream canceled")
			return
		case <-ticker.C:
			log.WithFields(log.Fields{
				"target":         name,
				"retry-counter":  counter,
				"offset-seconds": mod,
			}).Error("retry target setup")
			err := r.initTarget(ctx, name)
			if err == nil {
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

// TODO: check if status set correctly in the code!
func (r *Register) setTargetStatus(name string, status core.TargetStorageStatus) {
	r.Lock()
	defer r.Unlock()
	r.registerStatus[name] = status
}

func (r *Register) addTarget(name string, ts core.TargetStorage) {
	log.WithFields(log.Fields{
		"target": name,
	}).Info("target successfully added to the register")

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
