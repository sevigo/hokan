package directories

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/nicksnyder/basen"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
)

func HandleCreate(dirStore core.DirectoryStore, event core.EventCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := new(core.Directory)
		err := json.NewDecoder(r.Body).Decode(dir)
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: cannot unmarshal request body")
			handler.JSON_400(w, r, "invalid request body")
			return
		}

		// create new ID with base62 value of the path, don't allow ID for the request
		dir.ID = basen.Base62Encoding.EncodeToString([]byte(dir.Path))
		err = dirStore.Create(r.Context(), dir)
		if err != nil {
			handler.JSON_500(w, r, "can't store a new directory")
			return
		}

		err = event.Publish(r.Context(), &core.EventData{
			Type: core.WatchDirStart,
			Data: dir,
		})
		if err != nil {
			logger.FromRequest(r).WithError(err).Error("api: cannot store a new directory")
			handler.JSON_500(w, r, "can't store a new directory")
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, dir)
	}
}
