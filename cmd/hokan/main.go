package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/server"
)

func main() {
	config, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	initLogging(config)
	ctx := context.Background()

	app, err := InitializeApplication(config)
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: cannot initialize server")
	}

	g := errgroup.Group{}
	g.Go(func() error {
		logrus.WithFields(
			logrus.Fields{
				"host": config.Server.Host,
				"port": config.Server.Port,
				"url":  config.Server.Addr,
			},
		).Infoln("starting the http server")
		return app.server.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Fatalln("program terminated")
	}
}

func initLogging(c config.Config) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: c.Logging.Pretty,
	})
}

// application is the main struct.
type application struct {
	server *server.Server
}

func newApplication(server *server.Server) application {
	return application{
		server: server,
	}
}
