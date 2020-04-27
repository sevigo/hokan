package targets

import (
	"net/http"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
)

func HandleList(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targets := targets.AllConfigs()

		renderData := &core.TargetConfigsListResp{
			Targets: targets,
			Links:   createLinks(r),
			Meta: core.MetaDataResp{
				TotalItems: len(targets),
			},
		}
		handler.JSON_200(w, r, renderData)
	}
}

func createLinks(r *http.Request) []core.LinksResp {
	return []core.LinksResp{
		{
			Rel:    "self",
			Href:   r.URL.EscapedPath(),
			Method: r.Method,
		},
		{
			Rel:    "settings",
			Href:   r.URL.EscapedPath() + "/{target}/settings",
			Method: "POST",
		},
	}
}
