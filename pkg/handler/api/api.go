package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	dirs "github.com/sevigo/hokan/pkg/handler/api/directories"
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
	Dirs   core.DirectoryStore
	Files  core.FileStore
	Events core.EventCreator
}

func New(fileStore core.FileStore, dirStore core.DirectoryStore, events core.EventCreator) *Server {
	return &Server{
		Dirs:   dirStore,
		Files:  fileStore,
		Events: events,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// r.Use(logger.Middleware(s.Logger))

	cors := cors.New(corsOpts)
	r.Use(cors.Handler)

	r.Route("/directories", func(r chi.Router) {
		r.Post("/", dirs.HandleCreate(s.Dirs, s.Events))
		r.Get("/", dirs.HandleList(s.Dirs))
		r.Get("/{pathID}", dirs.HandleFind(s.Dirs))
	})

	// List all avaible endpoints
	r.Get("/", HandleAPIList())

	return r
}

func HandleAPIList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links := []core.LinksResp{
			{
				Rel:    "self",
				Href:   r.URL.EscapedPath(),
				Method: "GET",
			},
			{
				Rel:    "localDirs",
				Href:   "/api/directories",
				Method: "GET",
			},
			{
				Rel:    "version",
				Href:   "/version",
				Method: "GET",
			},
			{
				Rel:    "health",
				Href:   "/healthz",
				Method: "GET",
			},
		}
		renderData := &core.APIListResp{
			Links: links,
		}
		handler.JSON_200(w, r, renderData)
	}
}
