package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	// r.Use(logger.Middleware(s.Logger))

	r.Get("/version", HandleVersion)

	return r
}
