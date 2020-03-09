package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const ErrOnShutdownCode = 1

// A Server defines parameters for running an HTTP server.
type Server struct {
	Addr    string
	Cert    string
	Host    string
	Handler http.Handler
}

func (s Server) ListenAndServe(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}

	// Setup our Ctrl+C handler
	g.Go(func() error {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		killSignal := <-interrupt
		switch killSignal {
		case os.Interrupt:
			log.Print("Got SIGINT...")
		case syscall.SIGTERM:
			log.Print("Got SIGTERM...")
		}
		err := s1.Shutdown(ctx)
		if err != nil {
			log.Error("Shutdown error")
			os.Exit(ErrOnShutdownCode)
		} else {
			os.Exit(0)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		return s1.Shutdown(ctx)
	})

	g.Go(s1.ListenAndServe)

	return g.Wait()
}
