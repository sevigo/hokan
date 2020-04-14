package files

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleList(fileStore core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")
		l := logger.FromRequest(r).WithField("target", targetName)
		l.Infof("files.HandleList(): q=%#v\n", r.URL.Query())

		data, err := fileStore.List(r.Context(), &core.FileListOptions{
			TargetName: targetName,
		})
		if errors.Is(err, core.ErrTargetNotActive) {
			l.WithError(err).Error("api: can't get target")
			handler.JSON_404(w, r, core.ErrorResp{Code: http.StatusNotFound, Msg: "target not found"})
			return
		}
		if err != nil {
			l.WithError(err).Error("api: can't list files")
			handler.JSON_400(w, r, core.ErrorResp{Code: http.StatusBadRequest, Msg: err.Error()})
			return
		}

		renderData := &core.FilesListResp{
			Files: data,
			Links: createLinks(r),
			Meta:  core.MetaDataResp{},
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
	}
}
