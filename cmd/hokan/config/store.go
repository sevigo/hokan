package config

import (
	"os"
	"path/filepath"

	home "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
)

func defaultStore(c *Config) {
	if c.Database.Path == "" {
		homeDir, _ := home.Dir()
		appDir := filepath.Join(homeDir, ".hokan")
		createDirIfNotExist(appDir)
		c.Database.Path = filepath.Join(appDir, "store.db")
	}
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot create application directory")
		}
	}
}
