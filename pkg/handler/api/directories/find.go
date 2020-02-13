package directories

import (
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api/render"
)

func HandleFind(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, err := base64.StdEncoding.DecodeString(chi.URLParam(r, "path"))
		if err != nil {
			render.ErrorCode(w, err, 500)
			// logger.FromRequest(r).WithError(err).
			// 	Errorln("api: invlid directory")
			return
		}

		dir, err := dirStore.FindName(r.Context(), string(path))
		if err != nil {
			render.ErrorCode(w, err, 400)
			// logger.FromRequest(r).WithError(err).
			// 	Errorln("api: invlid directory")
			return
		}

		render.JSON(w, dir, 200)
	}
}
