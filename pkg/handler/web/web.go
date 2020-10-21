package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

var corsOpts = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

type Server struct {
	Logger *logrus.Logger
	SSE    core.ServerSideEventCreator
}

func New(logger *logrus.Logger, sse core.ServerSideEventCreator) *Server {
	return &Server{
		Logger: logger,
		SSE:    sse,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)

	cors := cors.New(corsOpts)
	r.Use(cors.Handler)

	s.SSE.PublishMessage("ping")

	r.Get("/version", HandleVersion)
	r.Get("/info", HandleInfo)
	r.Get("/events", s.SSE.Handler)

	return r
}
