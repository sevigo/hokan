package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/server"
)

func main() {
	config, err := config.Environ()
	if err != nil {
		log.Fatal().Err(err).Msg("main: invalid configuration")
	}

	initLogging(config)
	ctx := context.Background()

	app, err := InitializeApplication(config)
	if err != nil {
		log.Fatal().Err(err).Msg("main: cannot initialize server")
	}

	g := errgroup.Group{}
	g.Go(func() error {
		log.Info().
			Str("port", config.Server.Port).
			Str("url", config.Server.Addr).
			Msg("starting the http server")
		return app.server.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("main: program terminated")
	}
}

func initLogging(c config.Config) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822})
}

// application is the main struct.
type application struct {
	dirs   core.DirectoryStore
	server *server.Server
}

func newApplication(server *server.Server, dirs core.DirectoryStore) application {
	return application{
		dirs:   dirs,
		server: server,
	}
}
