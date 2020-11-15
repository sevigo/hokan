package main

import (
	"github.com/google/wire"
	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/gui"
)

var guiSet = wire.NewSet(
	provideAppConfig,
	provideGUIServer,
)

func provideGUIServer(guiConfig *core.AppConfig) *gui.Server {
	return gui.NewServer(guiConfig)
}

func provideAppConfig(config config.Config) *core.AppConfig {
	return &core.AppConfig{
		AppName:      config.GUI.AppName,
		BuildDir:     config.GUI.BuildDir,
		ResourcesDir: config.GUI.ResourcesDir,
		DevTools:     config.GUI.DevTools,
	}
}
