package directories

import (
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
			render.Status(r, 401)
			logger.FromRequest(r).Err(err).Msg("api: cannot unmarshal request body")
			return
		}

		err = dir.Validate()
		if err != nil {
			render.Status(r, 400)
			logger.FromRequest(r).Err(err).Msg("api: invlid directory")
			return
		}

		err = dirStore.Create(r.Context(), dir)
		if err != nil {
			render.Status(r, 500)
			logger.FromRequest(r).Err(err).Msg("api: cannot store a new directory")
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, dir)
	}
}
