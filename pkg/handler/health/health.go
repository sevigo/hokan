package health

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "OK")
	}
}
