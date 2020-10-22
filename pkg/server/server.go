package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const ErrOnShutdownCode = 1
const hardTimeout = 5 * time.Second

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

		// using SSE stream not allways allows to shutdown the service gracefully,
		// so we use a hard timeout here
		newCtx, cancel := context.WithTimeout(context.Background(), hardTimeout)
		defer cancel()

		if err := s1.Shutdown(newCtx); err != nil {
			log.Errorf("Shutdown error: %v", err)
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
