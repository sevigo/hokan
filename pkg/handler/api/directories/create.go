package directories

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleCreate(dirStore core.DirectoryStore, event core.EventCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := new(core.Directory)
		err := json.NewDecoder(r.Body).Decode(dir)
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: cannot unmarshal request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusBadRequest, Msg: "invalid request body"})
			return
		}

		err = dirStore.Create(r.Context(), dir)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusInternalServerError, Msg: "cannot store a new directory"})
			logger.FromRequest(r).WithError(err).Error("api: cannot store a new directory")
			return
		}

		dir.ID = base64.StdEncoding.EncodeToString([]byte(dir.Path))
		err = event.Publish(r.Context(), &core.EventData{
			Type: core.WatchDirStart,
			Data: dir,
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, core.ErrorResp{Code: http.StatusInternalServerError, Msg: "cannot store a new directory"})
			logger.FromRequest(r).WithError(err).Error("api: cannot store a new directory")
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, dir)
	}
}
