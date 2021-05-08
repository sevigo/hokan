package config

import (
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func defaultViodBackup() Backup {
	return Backup{
		Name: "void",
	}
}

func defaultLocalBackup() Backup {
	var localPath string
	if runtime.GOOS == "windows" {
		localPath = path.Join("c:", "backup")
	} else {
		localPath = path.Join("~", "backup")
	}

	return Backup{
		Name:            "local",
		TargetLocalPath: localPath,
	}
}

func defaultBackup(c *Config) {
	c.Backup = defaultViodBackup()
}

// TODO: we will read this config later from the UI maybe?
func defaultMinIO(c *Config) {
	viper.SetConfigType("yml")
	if c.Backup.configName == "" {
		c.Backup.configName = "config"
	}
	viper.SetConfigName(c.Backup.configName)
	if c.Backup.configPath == "" {
		c.Backup.configPath = "."
	}
	viper.AddConfigPath(c.Backup.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("Error reading config file")
	}
	err = viper.Unmarshal(&c.Backup)
	if err != nil {
		log.WithError(err).Fatal("Error unmarshaling config file")
	}
}
