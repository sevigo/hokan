package targets

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
)

func HandleList(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targets := targets.AllConfigs()

		renderData := map[string]interface{}{
			"targets": targets,
			"links": []core.LinksResp{
				{
					Rel:    "self",
					Href:   r.URL.EscapedPath(),
					Method: r.Method,
				},
				{
					Rel:    "files",
					Href:   r.URL.EscapedPath() + "/{targetName}/files",
					Method: r.Method,
				},
			},
		}
		render.Status(r, 200)
		render.JSON(w, r, renderData)
	}
}
