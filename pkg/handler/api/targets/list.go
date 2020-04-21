package targets

import (
	"net/http"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
)

func HandleList(targets core.TargetRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// targets := targets.AllTargets()

		renderData := &core.TargetListResp{
			Targets: map[string]core.Target{},
			Links:   createLinks(r),
			Meta: core.MetaDataResp{
				TotalItems: 0,
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
			Rel:    "files",
			Href:   r.URL.EscapedPath() + "/{targetName}/files",
			Method: r.Method,
		},
	}
}
