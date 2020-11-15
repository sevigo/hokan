package config

import (
	"os"
	"path"
)

func defaultGUI(c *Config) {
	if c.GUI.AppName == "" {
		c.GUI.AppName = "hokan"
	}
	c.GUI.DevTools = false
	pwd, _ := os.Getwd()
	c.GUI.ResourcesDir = path.Join(pwd, "resources")
	c.GUI.BuildDir = os.TempDir()
}
