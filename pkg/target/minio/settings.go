package minio

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

type Configurator struct {
	core.TargetFactory
}

func NewConfigurator() *Configurator {
	return &Configurator{New}
}

func (c *Configurator) Name() string {
	return TargetName
}

func (c *Configurator) Target() core.TargetFactory {
	return c.TargetFactory
}

func (c *Configurator) DefaultConfig() *core.TargetConfig {
	return &core.TargetConfig{
		Active:      false,
		Name:        "minio",
		Description: "open source cloud object storage server compatible with Amazon S3",
		Settings: core.TargetSettings{
			"MINIO_HOST":        "",
			"MINIO_ACCESS_KEY":  "",
			"MINIO_SECRET_KEY":  "",
			"MINIO_USE_SSL":     "",
			"MINIO_BUCKET_NAME": "",
		},
	}
}

func (c *Configurator) ValidateSettings(settings core.TargetSettings) (bool, error) {
	logger := log.WithField("target", TargetName)
	logger.Infof("ValidateSettings(): %+v", settings)

	for name := range c.DefaultConfig().Settings {
		value, ok := settings[name]
		if !ok {
			return false, fmt.Errorf("%q key is missing", name)
		}
		if value == "" {
			return false, fmt.Errorf("%q value is mepty", name)
		}
	}

	_, err := strconv.ParseBool(settings["MINIO_USE_SSL"])
	if err != nil {
		return false, fmt.Errorf("can't convert the value of MINIO_USE_SSL=%q to bool", settings["MINIO_USE_SSL"])
	}

	bucket := settings["MINIO_BUCKET_NAME"]
	match := bucketNameRegexp.MatchString(bucket)
	if !match {
		return false, fmt.Errorf("bucket name contains illegal characters")
	}

	return true, nil
}
