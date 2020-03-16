package targets

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleGet(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")

		conf, err := targets.GetConfig(r.Context(), targetName)
		if err != nil {
			render.Status(r, 400)
			logger.FromRequest(r).WithField("target", targetName).WithError(err).Error("api: can't get config")
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, conf)
	}
}
