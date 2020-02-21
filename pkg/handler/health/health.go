package health

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// New returns a new health check router
func New() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Handle("/", Handler())
	return r
}

// Handler creates an http.HandlerFunc that performs system healthchecks
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, 200)
		render.PlainText(w, r, "OK")
	}
}
