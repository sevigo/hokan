package config

import (
	"path"
	"runtime"
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
