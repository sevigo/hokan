package events

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sevigo/hokan/pkg/core"
)

func New(serverEvents core.ServerSideEventCreator) *Server {
	return &Server{
		serverEvents: serverEvents,
	}
}

type Server struct {
	serverEvents core.ServerSideEventCreator
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/events", s.serverEvents.Handler)

	return r
}
