package gui

import (
	"context"
	"fmt"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	log "github.com/sirupsen/logrus"
)

const windowHeight = 700
const windowWidth = 700

type Config struct {
	AppName      string
	BaseDir      string
	ResourcesDir string
}

func (c Config) Run(ctx context.Context) error {
	l := log.New()
	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           c.AppName,
		BaseDirectoryPath: c.BaseDir,
	})
	if err != nil {
		return fmt.Errorf("gui: creating astilectron failed: %w", err)
	}
	defer a.Close()
	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		return fmt.Errorf("gui: starting astilectron failed: %w", err)
	}

	// New window
	var useDevTools = true
	var w *astilectron.Window
	if w, err = a.NewWindow(c.BaseDir+"index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(windowHeight),
		Width:  astikit.IntPtr(windowWidth),
		WebPreferences: &astilectron.WebPreferences{
			DevTools: &useDevTools,
		},
	}); err != nil {
		return fmt.Errorf("gui: new window failed: %w", err)
	}

	// Create windows
	if err = w.Create(); err != nil {
		return fmt.Errorf("gui: creating window failed: %w", err)
	}

	// Blocking pattern
	a.Wait()
	return nil
}
