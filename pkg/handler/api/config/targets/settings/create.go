package settings

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

// POST /api/config/targets/{target}/settings
func HandleCleate(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "target")
		l := logger.FromRequest(r).WithField("target", targetName)

		configurator := targets.GetTargetStorageConfigurator(targetName)
		if configurator == nil {
			l.Error("api: can't find target")
			handler.JSON_404(w, r, "target not found")
			return
		}
		settings := core.TargetSettings{}
		err := json.NewDecoder(r.Body).Decode(&settings)
		if err != nil {
			l.WithError(err).Error("api: cannot unmarshal request body")
			handler.JSON_400(w, r, "invalid request body")
			return
		}
		if ok, err := configurator.ValidateSettings(settings); !ok {
			l.WithError(err).Error("api: invalid settings")
			handler.JSON_400(w, r, "invalid settings")
			return
		}
		config, err := targets.GetConfig(r.Context(), targetName)
		if err != nil {
			l.WithError(err).Error("api: cannot find default config")
			handler.JSON_404(w, r, "config not found")
			return
		}

		config.Settings = settings
		config.Active = true
		saveErr := targets.SetConfig(r.Context(), config)
		if saveErr != nil {
			l.WithError(err).Error("api: cannot store a new config")
			handler.JSON_404(w, r, "cannot store a new config")
			return
		}

		l.Info("target storage config saved successfully")
		handler.JSON_201(w, r, "target storage config saved successfully")
	}
}
