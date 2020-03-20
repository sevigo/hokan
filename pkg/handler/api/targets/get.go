package targets

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleGet(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")
		l := logger.FromRequest(r).WithField("target", targetName)
		l.Info("targets.HandleGet()")

		conf, err := targets.GetConfig(r.Context(), targetName)
		if errors.Is(err, core.ErrTargetConfigNotFound) {
			// TODO: combine this to one call!
			l.WithError(err).Error("api: can't get config")
			render.Status(r, 404)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusNotFound, Msg: err.Error()})
			return
		}
		if err != nil {
			l.WithError(err).Error("api: can't get config")
			render.Status(r, 500)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusInternalServerError, Msg: err.Error()})
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, conf)
	}
}
