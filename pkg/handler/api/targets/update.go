package targets

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleUpdate(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")

		_, err := targets.GetConfig(r.Context(), targetName)
		if err != nil {
			render.Status(r, 400)
			logger.FromRequest(r).WithField("target", targetName).WithError(err).Error("api: can't get config")
			return
		}

		conf := new(core.TargetConfig)
		err = json.NewDecoder(r.Body).Decode(conf)
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: cannot unmarshal request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusBadRequest, Msg: "invalid request body"})
			return
		}

		err = targets.SetConfig(r.Context(), conf)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusInternalServerError, Msg: "cannot store a new config"})
			logger.FromRequest(r).WithError(err).Error("api: cannot store a new config")
			return
		}

		logger.FromRequest(r).WithField("target", targetName).Infof("target config updated: %#v", conf)
		render.Status(r, http.StatusCreated)
	}
}
