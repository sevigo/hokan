package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	dirs "github.com/sevigo/hokan/pkg/handler/api/directories"
	targetstorage "github.com/sevigo/hokan/pkg/handler/api/targets"

	"github.com/sevigo/hokan/pkg/logger"
)

type Server struct {
	Logger  *logrus.Logger
	Dirs    core.DirectoryStore
	Events  core.EventCreator
	Targets core.TargetRegister
}

func New(dirStore core.DirectoryStore, events core.EventCreator, targets core.TargetRegister, logger *logrus.Logger) *Server {
	return &Server{
		Logger:  logger,
		Dirs:    dirStore,
		Events:  events,
		Targets: targets,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Use(logger.Middleware(s.Logger))

	r.Route("/directories", func(r chi.Router) {
		r.Post("/", dirs.HandleCreate(s.Dirs, s.Events))
		r.Get("/", dirs.HandleList(s.Dirs))
		r.Get("/{path}", dirs.HandleFind(s.Dirs))
	})

	r.Route("/targets", func(r chi.Router) {
		r.Get("/", targetstorage.HandleList(s.Targets))
	})

	return r
}
