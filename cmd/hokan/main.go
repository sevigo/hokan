package main

import (
	"context"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/server"
)

func main() {
	conf, err := config.Environ()
	if err != nil {
		log.Fatal().Err(err).Msg("main: invalid configuration")
	}

	initLogging()
	ctx := context.Background()

	app, err := InitializeApplication(ctx, conf)
	if err != nil {
		log.Fatal().Err(err).Msg("main: cannot initialize server")
	}

	g := errgroup.Group{}
	g.Go(func() error {
		log.Info().
			Str("port", conf.Server.Port).
			Str("url", conf.Server.Addr).
			Msg("starting the http server")
		return app.server.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("main: program terminated")
	}
}

func initLogging() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC822})
}

// application is the main struct.
type application struct {
	dirs    core.DirectoryStore
	watcher core.Watcher
	server  *server.Server
}

func newApplication(srv *server.Server, dirs core.DirectoryStore, watcher core.Watcher) application {
	return application{
		dirs:    dirs,
		server:  srv,
		watcher: watcher,
	}
}
