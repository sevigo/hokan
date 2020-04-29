package targets

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

// POST /api/config/targets/{target}/activate
func HandleActivate(targets core.TargetRegister) http.HandlerFunc {
	return triggerActivation(targets, true)
}

// POST /api/config/targets/{target}/deactivate
func HandleDeactivate(targets core.TargetRegister) http.HandlerFunc {
	return triggerActivation(targets, false)
}

// func RequestLogger(f LogFormatter) func(next http.Handler) http.Handler {
func triggerActivation(targets core.TargetRegister, active bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "target")
		log := logger.FromRequest(r).WithField("target", targetName)

		config, err := targets.GetConfig(r.Context(), targetName)
		if err != nil {
			log.WithError(err).Error("api: cannot find default config")
			handler.JSON_404(w, r, core.ErrorResp{Code: http.StatusNotFound, Msg: "config not found"})
			return
		}

		config.Active = active
		saveErr := targets.SetConfig(r.Context(), config)
		if saveErr != nil {
			log.WithError(err).Error("api: cannot store a new config")
			handler.JSON_400(w, r, core.ErrorResp{Code: http.StatusInternalServerError, Msg: "cannot store a new config"})
			return
		}

		log.Info("target activated successfully")
		handler.JSON_201(w, r)
	}
}
