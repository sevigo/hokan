package files

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleList(fileStore core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")
		l := logger.FromRequest(r).WithField("target", targetName)
		l.Infof("files.HandleList(): q=%#v\n", r.URL.Query())
		//  param1 := r.URL.Query().Get("param1")

		data, err := fileStore.List(r.Context(), &core.FileListOptions{
			TargetName: targetName,
		})
		if errors.Is(err, core.ErrTargetNotActive) {
			l.WithError(err).Error("api: can't get target")
			render.Status(r, 404)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusNotFound, Msg: err.Error()})
			return
		}
		if err != nil {
			l.WithError(err).Error("api: can't list files")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusBadRequest, Msg: err.Error()})
			return
		}

		renderData := map[string]interface{}{
			"files": data,
			"links": []core.LinksResp{
				{
					Rel:    "self",
					Href:   r.URL.EscapedPath(),
					Method: r.Method,
				},
			},
		}
		render.Status(r, 200)
		render.JSON(w, r, renderData)
	}
}
