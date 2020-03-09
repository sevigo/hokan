package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	dirs "github.com/sevigo/hokan/pkg/handler/api/directories"
)

type Server struct {
	Logger logrus.Logger
	Dirs   core.DirectoryStore
	Events core.EventCreator
}

func New(dirStore core.DirectoryStore, events core.EventCreator) Server {
	return Server{
		Dirs:   dirStore,
		Events: events,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// r.Use(logger.Middleware())

	r.Route("/directories", func(r chi.Router) {
		r.Post("/", dirs.HandleCreate(s.Dirs, s.Events))
		r.Get("/", dirs.HandleList(s.Dirs))
		r.Get("/{path}", dirs.HandleFind(s.Dirs))
	})

	return r
}
