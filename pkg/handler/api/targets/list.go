package targets

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
)

func HandleList(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targets := targets.AllConfigs()

		render.Status(r, 200)
		render.JSON(w, r, targets)
	}
}
