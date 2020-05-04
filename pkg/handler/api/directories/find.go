package directories

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/directory"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleFind(dirStore core.DirectoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathID := chi.URLParam(r, "pathID")
		dir, err := dirStore.FindName(r.Context(), pathID)
		if errors.Is(err, core.ErrDirectoryNotFound) {
			logger.FromRequest(r).WithError(err).Error("api: invalid directory")
			handler.JSON_404(w, r, "invalid directory name")
			return
		}
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: invalid directory")
			handler.JSON_500(w, r, "invalid directory")
			return
		}

		stats, err := directory.Stats(dir.Path)
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: invalid directory")
			handler.JSON_500(w, r, "invalid directory")
			return
		}

		renderData := core.DirectoryDetails{
			Directory: *dir,
			Stats:     *stats,
			Links:     createDirectoryDetailsLinks(r),
		}

		handler.JSON_200(w, r, renderData)
	}
}

func createDirectoryDetailsLinks(r *http.Request) []core.LinksResp {
	return []core.LinksResp{
		{
			Rel:    "self",
			Href:   r.URL.EscapedPath(),
			Method: r.Method,
		},
	}
}
