package directories

import (
	"net/http"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api/render"
)

func HandleList(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dirs, err := dirStore.List(r.Context())
		if err != nil {
			render.InternalError(w, err)
			// logger.FromRequest(r).WithError(err).
			// 	Warnln("api: cannot list users")
			return
		}

		render.JSON(w, dirs, 200)
	}
}
