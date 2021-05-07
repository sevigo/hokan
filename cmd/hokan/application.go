package main

import (
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/server"
)

// application is the main struct.
type application struct {
	dirs    core.DirectoryStore
	watcher core.Watcher
	backup  core.Backup
	server  *server.Server
}

func newApplication(srv *server.Server, dirs core.DirectoryStore, watcher core.Watcher, backup core.Backup) application {
	return application{
		dirs:    dirs,
		server:  srv,
		watcher: watcher,
		backup:  backup,
	}
}
