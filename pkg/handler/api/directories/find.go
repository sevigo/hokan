package directories

import (
	"encoding/base64"
	"net/http"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleFind(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, err := base64.StdEncoding.DecodeString(chi.URLParam(r, "path"))
		if err != nil {
			render.Status(r, 500)
			logger.FromRequest(r).Err(err).Msg("api: cannot encode path")
			return
		}

		dir, err := dirStore.FindName(r.Context(), string(path))
		if err != nil {
			render.Status(r, 400)
			logger.FromRequest(r).Err(err).Msg("api: invlid directory")
			return
		}
		render.Status(r, 200)
		render.JSON(w, r, dir)
	}
}
