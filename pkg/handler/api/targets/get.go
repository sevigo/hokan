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
			render.Status(r, 400)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusBadRequest, Msg: err.Error()})
			return
		}

		renderData := map[string]interface{}{
			"active":      conf.Active,
			"name":        conf.Name,
			"description": conf.Description,
			"settings":    conf.Settings,
			"links": []core.LinksResp{
				{
					Rel:    "self",
					Href:   r.URL.EscapedPath(),
					Method: r.Method,
				},
				{
					Rel:    "files",
					Href:   r.URL.EscapedPath() + "/files",
					Method: "GET",
				},
				{
					Rel:    "edit",
					Href:   r.URL.EscapedPath(),
					Method: "PUT",
				},
			},
		}
		render.Status(r, 200)
		render.JSON(w, r, renderData)
	}
}
