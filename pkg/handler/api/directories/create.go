package directories

import (
	"net/http"

	"github.com/sevigo/hokan/pkg/core"
)

func HandleCreate(dir core.DirectoryStore, service core.DirectoryService, event core.EventCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
