package server

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

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
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return s1.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		return s1.ListenAndServe()
	})
	return g.Wait()
}
