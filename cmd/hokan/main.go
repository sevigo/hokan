package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/gui"
	"github.com/sevigo/hokan/pkg/server"
)

func main() {
	conf, err := config.Environ()
	if err != nil {
		log.Fatal("main: invalid configuration")
	}

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

	g.Go(func() error {
		log.Info("starting GUI")
		return app.gui.Run(ctx)
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
	gui     *gui.Config
}

func newApplication(
	srv *server.Server,
	gui *gui.Config,
	dirs core.DirectoryStore,
	watcher core.Watcher,
	targets core.TargetRegister) application {
	return application{
		dirs:    dirs,
		server:  srv,
		watcher: watcher,
		targets: targets,
		gui:     gui,
	}
}
