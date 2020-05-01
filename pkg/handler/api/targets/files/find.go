package files

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleGet(fileStore core.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "targetName")
		fileID := chi.URLParam(r, "fileID")
		l := logger.FromRequest(r).WithFields(log.Fields{
			"target": targetName,
			"file":   fileID,
		})
		l.Infof("files.HandleGet(): q=%#v\n", r.URL.Query())

		data, err := fileStore.Find(r.Context(), &core.FileSearchOptions{
			ID:         fileID,
			TargetName: targetName,
		})
		if errors.Is(err, core.ErrFileNotFound) {
			l.WithError(err).Error("api: can't find file")
			handler.JSON_404(w, r, "can't find file")
			return
		}
		if err != nil {
			l.WithError(err).Error("api: can't get file")
			handler.JSON_500(w, r, "can't get file")
			return
		}

		handler.JSON_200(w, r, data)
	}
}
