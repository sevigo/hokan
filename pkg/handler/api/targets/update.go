package targets

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleUpdate(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")
		l := logger.FromRequest(r).WithField("target", targetName)

		oldConf, err := targets.GetConfig(r.Context(), targetName)
		if err != nil {
			l.WithError(err).Error("api: can't get config")
			handler.JSON_400(w, r, "can't get config")
			return
		}

		conf := new(core.TargetConfig)
		err = json.NewDecoder(r.Body).Decode(conf)
		if err != nil {
			l.WithError(err).Error("api: cannot unmarshal request body")
			handler.JSON_400(w, r, "invalid request body")
			return
		}

		// Description and Name are read-only, so we just overwrite these fields
		conf.Description = oldConf.Description
		conf.Name = oldConf.Name

		err = targets.SetConfig(r.Context(), conf)
		if err != nil {
			l.WithError(err).Error("api: can't store a new config")
			handler.JSON_500(w, r, "can't store a new config")
			return
		}

		l.Info("targets.HandleUpdate(): target storage config updated successfully")
		handler.JSON_201(w, r, "target storage config updated successfully")
	}
}
