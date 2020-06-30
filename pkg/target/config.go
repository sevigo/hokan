package target

import (
	"context"
	"errors"
	"fmt"

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

func (r *Register) GetConfig(ctx context.Context, name string) (*core.TargetConfig, error) {
	var err error
	target := r.GetTarget(name)
	fmt.Printf(">>>> GetConfig(): target=%v\n", target)
	// defaultConf := target.DefaultConfig()
	conf, err := r.configStore.Find(ctx, name)
	if errors.Is(err, configstore.ErrConfigNotFound) {
		err = r.configStore.Save(ctx, defaultConf)
		conf = defaultConf
	}
	return conf, err
}

func (r *Register) SetConfig(ctx context.Context, config *core.TargetConfig) error {
	err := r.configStore.Save(ctx, config)
	if err != nil {
		return err
	}
	log.WithField("target", config.Name).Info("target.SetConfig(): new config stored successfully")
	// we changed from passiv to active, so we init the storage and rescan
	if config.Active {
		err = r.initTarget(ctx, config.Name)
		if err != nil {
			// initWithRetry will call rescanAllWatchedDirs
			go r.initWithRetry(ctx, config.Name)
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
