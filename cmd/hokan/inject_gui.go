package main

import (
	"github.com/google/wire"
	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/gui"
)

var guiSet = wire.NewSet(
	provideGUI,
)

func provideGUI(config config.Config) *gui.Config {
	return &gui.Config{
		AppName: config.GUI.AppName,
		BaseDir: config.GUI.BaseDir,
	}
}
