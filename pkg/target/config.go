package target

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	configstore "github.com/sevigo/hokan/pkg/store/config"
)

func (r *Register) AllConfigs() map[string]core.TargetConfig {
	configs := map[string]core.TargetConfig{}

	for name := range r.AllTargets() {
		conf, err := r.GetConfig(r.ctx, name)
		if err != nil {
			log.WithError(err).Errorf("can't get config for %q", name)
			continue
		}
		configs[name] = *conf
	}

	return configs
}

func (r *Register) GetConfig(ctx context.Context, name string) (conf *core.TargetConfig, err error) {
	conf, err = r.configStore.Find(ctx, name)
	if errors.Is(err, configstore.ErrConfigNotFound) {
		defaultConf, ok := r.storagesDefaultConfigs[name]
		if !ok {
			return nil, configstore.ErrConfigNotFound
		}
		err = r.configStore.Save(ctx, defaultConf)
		conf = defaultConf
	}
	return
}

func (r *Register) SetConfig(ctx context.Context, config *core.TargetConfig) error {
	err := r.configStore.Save(ctx, config)
	if err != nil {
		return err
	}
	log.WithField("target", config.Name).Info("target.SetConfig(): new config stored successfully")
	// we changed from passiv to active, so we init the storage and rescan
	if config.Active {
		configurator, ok := r.configurators[config.Name]
		if !ok {
			return core.ErrTargetConfigNotFound
		}
		err = r.initTarget(ctx, configurator)
		if err != nil {
			// initWithRetry will call rescanAllWatchedDirs
			go r.initWithRetry(ctx, configurator)
		} else {
			log.WithField("target", config.Name).Info("target.SetConfig(): target is activated now")
			r.rescanAllWatchedDirs()
			r.setTargetStatus(config.Name, core.TargetStorageOK)
		}
	} else {
		log.WithField("target", config.Name).Info("target.SetConfig(): target is deactivated now")
		r.setTargetStatus(config.Name, core.TargetStoragePaused)
	}
	return nil
}
