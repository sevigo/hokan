package files

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sevigo/hokan/pkg/core"
)

func HandleFind(fileStore core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, 200)
		render.PlainText(w, r, "")
	}
}
