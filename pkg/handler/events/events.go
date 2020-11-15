package events

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sevigo/hokan/pkg/core"
)

type Server struct {
	serverEvents core.ServerSideEventCreator
}

func New(serverEvents core.ServerSideEventCreator) *Server {
	return &Server{
		serverEvents: serverEvents,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/events", s.serverEvents.Handler)

	return r
}
