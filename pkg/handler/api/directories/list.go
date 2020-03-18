package directories

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleList(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dirs, err := dirStore.List(r.Context())
		if err != nil {
			render.Status(r, 500)
			logger.FromRequest(r).WithError(err).Error("api: cannot list directories")
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, dirs)
	}
}
