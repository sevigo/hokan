package target

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	configstore "github.com/sevigo/hokan/pkg/store/config"
	"github.com/sevigo/hokan/pkg/target/local"
	"github.com/sevigo/hokan/pkg/target/minio"
	"github.com/sevigo/hokan/pkg/target/void"
)

var defaultConfigs = map[string]*core.TargetConfig{
	local.TargetName: local.DefaultConfig(),
	minio.TargetName: minio.DefaultConfig(),
	void.TargetName:  void.DefaultConfig(),
}

func (r *Register) AllConfigs() map[string]*core.TargetConfig {
	return defaultConfigs
}

func (r *Register) GetConfig(ctx context.Context, targetName string) (*core.TargetConfig, error) {
	var err error
	defaultConf, ok := defaultConfigs[targetName]
	if !ok {
		log.Errorf("default config for target storage %q not found", targetName)
		return nil, core.ErrTargetConfigNotFound
	}
	conf, err := r.configStore.Find(ctx, targetName)
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
