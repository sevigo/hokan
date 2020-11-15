package tray

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

func GetMainMenu(a *astilectron.Astilectron) (*astilectron.Window, error) {
	w, err := a.NewWindow("http://localhost:8081/gui", &astilectron.WindowOptions{
		AlwaysOnTop:     astikit.BoolPtr(false),
		AutoHideMenuBar: astikit.BoolPtr(true),
		BackgroundColor: astikit.StrPtr("#2e2c29"),
		Center:          astikit.BoolPtr(true),
		Closable:        astikit.BoolPtr(true),
		Height:          astikit.IntPtr(550),
		Resizable:       astikit.BoolPtr(true),
		Title:           astikit.StrPtr("Hokan"),
		TitleBarStyle:   astikit.StrPtr("hidden"),
		Width:           astikit.IntPtr(300),
	})
	if err != nil {
		return nil, err
	}

	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)
		// Process message
		if s == "hello" {
			return "world"
		}
		return nil
	})

	return w, nil
}
