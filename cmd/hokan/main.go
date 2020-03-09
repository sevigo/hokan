package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/server"
)

func main() {
	conf, err := config.Environ()
	if err != nil {
		log.Fatal("main: invalid configuration")
	}

	initLogging()
	ctx := context.Background()

	app, err := InitializeApplication(ctx, conf)
	if err != nil {
		log.Fatal("main: cannot initialize server")
	}

	g := errgroup.Group{}
	g.Go(func() error {
		log.WithFields(log.Fields{
			"port": conf.Server.Port,
		}).Info("starting the http server")
		return app.server.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal("main: program terminated")
	}
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.Info("Hello Hokan!")
}

// application is the main struct.
type application struct {
	dirs    core.DirectoryStore
	watcher core.Watcher
	target  core.Target
	server  *server.Server
}

func newApplication(srv *server.Server, dirs core.DirectoryStore, watcher core.Watcher, target core.Target) application {
	return application{
		dirs:    dirs,
		server:  srv,
		watcher: watcher,
		target:  target,
	}
}
