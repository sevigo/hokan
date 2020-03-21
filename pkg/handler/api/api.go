package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	dirs "github.com/sevigo/hokan/pkg/handler/api/directories"
	targetstorage "github.com/sevigo/hokan/pkg/handler/api/targets"
	targetsfiles "github.com/sevigo/hokan/pkg/handler/api/targets/files"

	"github.com/sevigo/hokan/pkg/logger"
)

type Server struct {
	Logger  *logrus.Logger
	Dirs    core.DirectoryStore
	Files   core.FileStore
	Events  core.EventCreator
	Targets core.TargetRegister
}

func New(fileStore core.FileStore, dirStore core.DirectoryStore, events core.EventCreator, targets core.TargetRegister, logger *logrus.Logger) *Server {
	return &Server{
		Logger:  logger,
		Dirs:    dirStore,
		Files:   fileStore,
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
		r.Get("/{pathID}", dirs.HandleFind(s.Dirs))
	})

	r.Route("/targets", func(r chi.Router) {
		r.Put("/{targetName}", targetstorage.HandleUpdate(s.Targets))
		r.Get("/", targetstorage.HandleList(s.Targets))
		r.Get("/{targetName}", targetstorage.HandleGet(s.Targets))

		r.Route("/{targetName}/files", func(r chi.Router) {
			r.Get("/", targetsfiles.HandleList(s.Files))
			r.Get("/{fileID}", targetsfiles.HandleGet(s.Files))
		})
	})

	return r
}
