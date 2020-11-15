package gui

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/sevigo/hokan/pkg/gui/tray"
)

const (
	electronV10    = "10.1.4"
	windowHeight   = 700
	windowWidth    = 700
	debugEventsURL = "http://localhost:8081/gui/debug"
)

// https://github.com/snight1983/BitcoinVIP/blob/582cc9975157c4f2517ab59966f5e656ebdf9ee3/spvwallet/gui/bootstrap/run.go
func (s *Server) Run(ctx context.Context) error {
	l := log.New(os.Stderr, "", 0)
	s.config.IconPath = path.Join(s.config.ResourcesDir, "icons", "icon-color.png")
	var a *astilectron.Astilectron
	var err error

	if a, err = astilectron.New(l, astilectron.Options{
		AppIconDefaultPath: s.config.IconPath,
		AppName:            s.config.AppName,
		// where you want the provisioner to install the dependencies
		BaseDirectoryPath: s.config.BuildDir,
		VersionElectron:   electronV10,
	}); err != nil {
		return fmt.Errorf("gui: creating astilectron failed: %w", err)
	}
	defer a.Close()
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		return fmt.Errorf("gui: starting astilectron failed: %w", err)
	}

	// if err := c.addDebugWindow(a); err != nil {
	// 	return err
	// }
	if err := s.addTray(a); err != nil {
		return err
	}
	if m, err = tray.GetMainMenu(a); err != nil {
		return err
	}

	// Blocking pattern
	a.Wait()
	return nil
}

// Docs: https://github.com/asticode/go-astilectron
func (s Server) addDebugWindow(a *astilectron.Astilectron) error {
	var w *astilectron.Window
	var err error

	if w, err = a.NewWindow(debugEventsURL, &astilectron.WindowOptions{
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
	if s.config.DevTools {
		if err = w.OpenDevTools(); err != nil {
			return fmt.Errorf("gui: opening dev tools faild: %w", err)
		}
	}
	return nil
}

var m *astilectron.Window
var once sync.Once

func (s Server) addTray(a *astilectron.Astilectron) error {
	t := a.NewTray(&astilectron.TrayOptions{
		Image:   astikit.StrPtr(s.config.IconPath),
		Tooltip: astikit.StrPtr(s.config.AppName),
	})

	// Create tray
	if err := t.Create(); err != nil {
		return err
	}

	t.On(astilectron.EventNameTrayEventClicked, func(e astilectron.Event) (deleteListener bool) {
		once.Do(func() {
			m.Create()
			if s.config.DevTools {
				m.OpenDevTools()
			}
		})
		if !m.IsShown() {
			m.Show()
		} else {
			m.Hide()
		}
		return
	})

	return nil
}
