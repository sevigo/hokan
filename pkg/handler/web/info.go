package web

import (
	"net/http"
	"os"
	"os/user"
	"runtime"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler"
	"github.com/sevigo/hokan/pkg/logger"
	"github.com/sevigo/hokan/pkg/version"
)

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	logger.FromRequest(r).Info("web.HandleInfo()")

	machine, err := os.Hostname()
	if err != nil {
		machine = "machine"
	}

	userName := "user"
	u, err := user.Current()
	if err == nil {
		userName = u.Name
	}

	info := &core.Info{
		Machine: machine,
		OS:      runtime.GOOS,
		User:    userName,
		Version: version.Version.String(),
	}

	handler.JSON_200(w, r, info)
}
