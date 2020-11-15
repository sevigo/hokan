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

// All known target configurations
var targetConfigurators = []core.TargetStorageConfigurator{
	local.NewConfigurator(),
	minio.NewConfigurator(),
	void.NewConfigurator(),
}

type Register struct {
	sync.Mutex

	ctx context.Context
	// write results of any operation here, this will be propagated to the UI
	Results     chan core.TargetOperationResult
	fileStore   core.FileStore
	configStore core.ConfigStore
	event       core.EventCreator

	storages               map[string]core.TargetStorage
	storagesStatus         map[string]core.TargetStorageStatus
	storagesDefaultConfigs map[string]*core.TargetConfig
	configurators          map[string]core.TargetStorageConfigurator
}

func New(
	ctx context.Context,
	fileStore core.FileStore,
	configStore core.ConfigStore,
	event core.EventCreator) (core.TargetRegister, error) {
	log.Debug("target.New(): start")

	r := &Register{
		ctx:         ctx,
		fileStore:   fileStore,
		configStore: configStore,
		event:       event,

		storages:               make(map[string]core.TargetStorage),
		storagesStatus:         make(map[string]core.TargetStorageStatus),
		storagesDefaultConfigs: make(map[string]*core.TargetConfig),
		configurators:          make(map[string]core.TargetStorageConfigurator),

		Results: make(chan core.TargetOperationResult),
	}

	r.initTargets(ctx)
	go r.StartFileAddedWatcher()
	return r, nil
}

func (r *Register) GetTargetStorageConfigurator(name string) core.TargetStorageConfigurator {
	return r.configurators[name]
}

func (r *Register) AllTargets() map[string]core.Target {
	result := map[string]core.Target{}
	for name, status := range r.storagesStatus {
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
	return r.storages[name]
}

func (r *Register) initTargets(ctx context.Context) {
	for _, configurator := range targetConfigurators {
		targetName := configurator.Name()
		r.configurators[targetName] = configurator
		r.storagesDefaultConfigs[targetName] = configurator.DefaultConfig()
		err := r.initTarget(ctx, configurator)
		if errors.Is(err, core.ErrTargetNotActive) {
			log.WithError(err).
				WithField("target", targetName).
				Error("initTargets(): ignore the target for now")
			continue
		}
		if err != nil {
			log.WithError(err).Error("initTargets()")
			go r.initWithRetry(ctx, configurator)
		}
	}
}

func (r *Register) initTarget(ctx context.Context, configurator core.TargetStorageConfigurator) error {
	name := configurator.Name()
	logger := log.WithFields(log.Fields{
		"target": name,
	})

	conf, err := r.GetConfig(ctx, name)
	if err != nil {
		logger.WithError(err).
			Fatal("can't get configuration for the target storage")
		return err
	}

	target := configurator.Target()
	ts, err := target(ctx, r.fileStore, *conf)
	if err != nil {
		return err
	}

	r.addTarget(name, ts)
	return nil
}

func (r *Register) initWithRetry(ctx context.Context, configurator core.TargetStorageConfigurator) {
	counter := 1
	mod := 10
	ticker := time.NewTicker(time.Duration(mod) * time.Second)
	defer ticker.Stop()

	name := configurator.Name()
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
			err := r.initTarget(ctx, configurator)
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
	r.storagesStatus[name] = status
}

func (r *Register) addTarget(name string, ts core.TargetStorage) {
	log.WithFields(log.Fields{
		"target": name,
	}).Info("target successfully added to the register")

	r.Lock()
	defer r.Unlock()
	r.storages[name] = ts
	r.storagesStatus[name] = core.TargetStorageOK
}

func (r *Register) getTargetStatus(name string) core.TargetStorageStatus {
	r.Lock()
	defer r.Unlock()
	return r.storagesStatus[name]
}
