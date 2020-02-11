package render

import (
	"encoding/json"
	"net/http"

	"github.com/sevigo/hokan/pkg/handler/api/errors"
)

func BadRequest(w http.ResponseWriter, err error) {
	ErrorCode(w, err, 400)
}

func InternalError(w http.ResponseWriter, err error) {
	ErrorCode(w, err, 500)
}

func ErrorCode(w http.ResponseWriter, err error, status int) {
	JSON(w, &errors.Error{Message: err.Error()}, status)
}

func JSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
