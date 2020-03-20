package directories

import (
	"errors"
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
		if errors.Is(err, core.ErrDirectoryNotFound) {
			logger.FromRequest(r).WithError(err).Error("api: invalid directory")
			render.Status(r, 404)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusNotFound, Msg: err.Error()})
			return
		}
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: invalid directory")
			render.Status(r, 500)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusInternalServerError, Msg: err.Error()})
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, dir)
	}
}
