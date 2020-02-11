package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/sevigo/hokan/pkg/core"
	dirs "github.com/sevigo/hokan/pkg/handler/api/directories"
)

type Server struct {
	Dirs   core.DirectoryStore
	Events core.EventCreator
}

func New(dirs core.DirectoryStore, events core.EventCreator) Server {
	return Server{
		Dirs:   dirs,
		Events: events,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/directories", func(r chi.Router) {
		r.Post("/", dirs.HandleCreate(s.Dirs, s.Events))
	})

	return r
}
