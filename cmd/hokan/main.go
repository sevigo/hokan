package main

import (
	"context"

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

	// initLogging()
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

// application is the main struct.
type application struct {
	dirs    core.DirectoryStore
	watcher core.Watcher
	targets core.TargetRegister
	server  *server.Server
}

func newApplication(srv *server.Server, dirs core.DirectoryStore, watcher core.Watcher, targets core.TargetRegister) application {
	return application{
		dirs:    dirs,
		server:  srv,
		watcher: watcher,
		targets: targets,
	}
}
