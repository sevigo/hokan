package directories

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleList(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dirs, err := dirStore.List(r.Context())
		if err != nil {
			render.Status(r, 500)
			logger.FromRequest(r).WithError(err).Error("api: cannot list directories")
			return
		}

		resp := &core.DirectoriesListResp{
			Directories: dirs,
			Links:       createLinks(r),
			Meta: core.MetaDataResp{
				TotalItems: len(dirs),
			},
		}
		handler.JSON_200(w, r, resp)
	}
}

func createLinks(r *http.Request) []core.LinksResp {
	return []core.LinksResp{
		{
			Rel:    "self",
			Href:   r.URL.EscapedPath(),
			Method: r.Method,
		},
	}
}
