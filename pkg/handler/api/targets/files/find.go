package files

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
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
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusNotFound, Msg: err.Error()})
			return
		}
		if err != nil {
			l.WithError(err).Error("api: can't get file")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusBadRequest, Msg: err.Error()})
			return
		}

		render.Status(r, 200)
		render.JSON(w, r, data)
	}
}

/*
{
	"ID": "1F8u2r1zHGKDjCSJ7myaswS5C97iDDhv7evEa74xj3xsuMCLg3OXItkgS6PiB",
	"Path":	"/home/igor/MyLocalFiles/pixiv/66093971_p0.jpg",
	"Checksum": "e50b49643d0ee2d19bf9adc9f02cd803a4ebe6ff333bb3ae2457a9f53820194f",
	"Info": "{\"ModTime\":\"2020-03-21T20:37:23.908618193+01:00\",\"Mode\":436,\"Name\":\"66093971_p0.jpg\",\"Size\":935942}",
	"Targets": ["void"],
	"links": [
		{
			"href":"/api/targets/void/files/1F8u2r1zHGKDjCSJ7myaswS5C97iDDhv7evEa74xj3xsuMCLg3OXItkgS6PiB",
			"rel":"self",
			"method":"GET",
		},
	]
}
*/
