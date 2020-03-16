package target

import (
	"context"
	"errors"
	"fmt"

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
		return nil, fmt.Errorf("default config for target storage %q not found", targetName)
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
	err = r.initTarget(ctx, config.Name)
	if err != nil {
		go r.initWithRetry(ctx, config.Name)
	} else {
		r.rescanAllWatchedDirs()
	}
	return nil
}
