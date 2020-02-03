package config

import "net/http"

// Config provides the system configuration.
type Config struct{}

type Server struct {
	Addr    string
	Host    string
	Port    string
	Handler http.Handler
}

func Environ() (Config, error) {
	cfg := Config{}
	return cfg, nil
}
