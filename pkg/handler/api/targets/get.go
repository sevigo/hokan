package targets

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleGet(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")
		l := logger.FromRequest(r).WithField("target", targetName)
		l.Info("targets.HandleGet()")

		conf, err := targets.GetConfig(r.Context(), targetName)
		if errors.Is(err, core.ErrTargetConfigNotFound) {
			l.WithError(err).Error("api: can't get config")
			handler.JSON_404(w, r, "can't get config")
			return
		}
		if err != nil {
			l.WithError(err).Error("api: can't get config")
			render.Status(r, 400)
			handler.JSON_400(w, r, "can't get config")
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

		handler.JSON_200(w, r, renderData)
	}
}
