package web

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sevigo/hokan/pkg/logger"
	"github.com/sevigo/hokan/pkg/version"
)

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	l := logger.FromRequest(r).WithField("version", version.Version.String())
	l.Info("web.HandleVersion()")

	v := struct {
		Source  string `json:"source,omitempty"`
		Version string `json:"version,omitempty"`
		Commit  string `json:"commit,omitempty"`
	}{
		Source:  version.GitRepository,
		Commit:  version.GitCommit,
		Version: version.Version.String(),
	}

	writeJSON(w, r, &v, 200)
}

func writeJSON(w http.ResponseWriter, r *http.Request, v interface{}, status int) {
	render.Status(r, status)
	render.JSON(w, r, v)
}
