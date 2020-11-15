package gui

import (
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/sevigo/hokan/pkg/core"
)

type Server struct {
	config *core.AppConfig
}

func NewServer(config *core.AppConfig) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		files := path.Join(s.config.ResourcesDir, "app")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(files)))
		fs.ServeHTTP(w, r)
	})

	return r
}
