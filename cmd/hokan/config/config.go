package config

import (
	"net/http"
)

// Config provides the system configuration.
type Config struct {
	Database Database
	Logging  Logging
	Server   Server
	GUI      GUI
}

type GUI struct {
	AppName      string
	DevTools     bool
	ResourcesDir string
	BuildDir     string
}

type Logging struct {
	Debug  bool
	Trace  bool
	Color  bool
	Pretty bool
	Text   bool
}

type Server struct {
	Addr    string
	Host    string
	Port    string
	Proto   string
	Handler http.Handler
}

type Database struct {
	Path string
}

func Environ() (Config, error) {
	cfg := Config{}
	defaultAddress(&cfg)
	defaultStore(&cfg)
	defaultLogger(&cfg)
	defaultGUI(&cfg)
	return cfg, nil
}

func defaultLogger(c *Config) {
	c.Logging = Logging{
		Debug: true,
		Color: true,
		Text:  true,
	}
}

func defaultAddress(c *Config) {
	if c.Server.Host != "" && c.Server.Proto != "" && c.Server.Port != "" {
		c.Server.Addr = c.Server.Proto + "://" + c.Server.Host + ":" + c.Server.Port
	} else {
		c.Server.Port = "8081"
		c.Server.Addr = ":" + c.Server.Port
	}
}
