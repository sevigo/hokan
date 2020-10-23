package gui

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	// log "github.com/sirupsen/logrus"
)

const windowHeight = 700
const windowWidth = 700
const entryPoint = "http://localhost:8081/debug/events/"

type Config struct {
	AppName      string
	BaseDir      string
	ResourcesDir string
	DevTools     bool
	IconPath     string
}

// https://github.com/snight1983/BitcoinVIP/blob/582cc9975157c4f2517ab59966f5e656ebdf9ee3/spvwallet/gui/bootstrap/run.go
func (c Config) Run(ctx context.Context) error {
	l := log.New(os.Stderr, "", 0)
	c.IconPath = path.Join(c.BaseDir, "icons", "icon-color.png")
	var a *astilectron.Astilectron
	var err error

	if a, err = astilectron.New(l, astilectron.Options{
		AppIconDefaultPath: c.IconPath,
		AppName:            c.AppName,
		// where you want the provisioner to install the dependencies
		BaseDirectoryPath: c.BaseDir,
		VersionElectron:   "10.1.4",
	}); err != nil {
		return fmt.Errorf("gui: creating astilectron failed: %w", err)
	}
	defer a.Close()
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		return fmt.Errorf("gui: starting astilectron failed: %w", err)
	}

	if err := c.addWindow(a); err != nil {
		return err
	}
	if err := c.addTray(a); err != nil {
		return err
	}

	// Blocking pattern
	a.Wait()
	return nil
}

// Docs: https://github.com/asticode/go-astilectron
func (c Config) addWindow(a *astilectron.Astilectron) error {
	var w *astilectron.Window
	var err error

	if w, err = a.NewWindow(entryPoint, &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(windowHeight),
		Width:  astikit.IntPtr(windowWidth),
	}); err != nil {
		return fmt.Errorf("gui: new window failed: %w", err)
	}

	// Create windows
	if err = w.Create(); err != nil {
		return fmt.Errorf("gui: creating window failed: %w", err)
	}
	// Open dev tools
	if c.DevTools {
		if err = w.OpenDevTools(); err != nil {
			return fmt.Errorf("gui: opening dev tools faild: %w", err)
		}
	}
	return nil
}

func (c Config) addTray(a *astilectron.Astilectron) error {
	var t = a.NewTray(&astilectron.TrayOptions{
		Image:   astikit.StrPtr(c.IconPath),
		Tooltip: astikit.StrPtr("Tray's tooltip"),
	})

	// Create tray
	if err := t.Create(); err != nil {
		return err
	}

	// New tray menu
	var m = t.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Root 1"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("Item 1")},
				{Label: astikit.StrPtr("Item 2")},
				{Type: astilectron.MenuItemTypeSeparator},
				{Label: astikit.StrPtr("Item 3")},
			},
		},
		{
			Label: astikit.StrPtr("Root 2"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("Item 1")},
				{Label: astikit.StrPtr("Item 2")},
			},
		},
	})

	// Create the menu
	return m.Create()
}
