package config

import (
	"net/http"
)

// Config provides the system configuration.
type Config struct {
	Database Database
	Logging  Logging
	Server   Server
	Backup   Backup
}

type Backup struct {
	configPath      string
	configName      string
	Name            string
	TargetLocalPath string
	MinIO           MinIOConf
}

// minio:
//     access_key_id: key
//     endpoint: localhost
//     secret_access_key: secret
//     use_ssl: false
type MinIOConf struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
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

// Environ returns a some default hard-coded configuration
// later it could be some configuration from the local config file
// or remote configuration service
func Environ() (Config, error) {
	cfg := Config{}
	defaultAddress(&cfg)
	defaultStore(&cfg)
	defaultLogger(&cfg)
	defaultBackup(&cfg)
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
