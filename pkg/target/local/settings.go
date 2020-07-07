package local

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sevigo/hokan/pkg/core"
	log "github.com/sirupsen/logrus"
)

const TargetName = "local"

var bucketNameRegexp = regexp.MustCompile("^[a-zA-Z0-9_.]+$")

var defaultConfig = core.TargetConfig{
	Active:      false,
	Name:        TargetName,
	Description: "store the files on the local disk",
	Settings: core.TargetSettings{
		"LOCAL_STORAGE_PATH": "",
		"LOCAL_BUCKET_NAME":  "",
	},
}

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
	return &defaultConfig
}

func (c *Configurator) ValidateSettings(settings core.TargetSettings) (bool, error) {
	logger := log.WithField("target", TargetName)
	logger.Infof("ValidateSettings(): %+v", settings)

	path, ok := settings["LOCAL_STORAGE_PATH"]
	if !ok {
		return false, fmt.Errorf("LOCAL_STORAGE_PATH is empty")
	}
	if _, err := os.Stat(filepath.Clean(path)); os.IsNotExist(err) {
		return false, fmt.Errorf("%q does not exist", filepath.Clean(path))
	}

	bucket, ok := settings["LOCAL_BUCKET_NAME"]
	if !ok {
		return false, fmt.Errorf("LOCAL_BUCKET_NAME is empty")
	}
	if bucket == "" {
		return false, fmt.Errorf("LOCAL_BUCKET_NAME is empty")
	}

	match := bucketNameRegexp.MatchString(bucket)
	if !match {
		return false, fmt.Errorf("bucket name contains illegal characters")
	}

	return true, nil
}
