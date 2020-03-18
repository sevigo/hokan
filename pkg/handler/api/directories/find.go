package directories

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleFind(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathID := chi.URLParam(r, "pathID")
		dir, err := dirStore.FindName(r.Context(), pathID)
		if err != nil {
			render.Status(r, 400)
			logger.FromRequest(r).WithError(err).Error("api: invalid directory")
			return
		}
		// TODO: add 404 check
		render.Status(r, 200)
		render.JSON(w, r, dir)
	}
}
