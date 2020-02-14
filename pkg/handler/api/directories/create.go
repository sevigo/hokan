package directories

import (
	"encoding/json"
	"net/http"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api/render"
)

func HandleCreate(dirStore core.DirectoryStore, event core.EventCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := new(core.Directory)
		err := json.NewDecoder(r.Body).Decode(dir)
		if err != nil {
			render.BadRequest(w, err)
			// logger.FromRequest(r).WithError(err).
			// 	Debugln("api: cannot unmarshal request body")
			return
		}

		err = dir.Validate()
		if err != nil {
			render.ErrorCode(w, err, 400)
			// logger.FromRequest(r).WithError(err).
			// 	Errorln("api: invlid directory")
			return
		}

		err = dirStore.Create(r.Context(), dir)
		if err != nil {
			render.InternalError(w, err)
			// logger.FromRequest(r).WithError(err).
			// 	Warnln("api: cannot store a new directory")
			return
		}

		render.JSON(w, dir, 200)
	}
}
