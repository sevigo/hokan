package main

import (
	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
	"github.com/sevigo/hokan/pkg/store/directory"
	"github.com/sevigo/hokan/pkg/store/file"
)

var storeSet = wire.NewSet(
	provideDatabase,
	provideDirectoryStore,
	provideFileStore,
)

func provideDatabase(config config.Config) (core.DB, error) {
	return db.Connect(config.Database.Path)
}

func provideDirectoryStore(db core.DB) core.DirectoryStore {
	dirs := directory.New(db)
	return dirs
}

func provideFileStore(db core.DB) core.FileStore {
	files := file.New(db)
	return files
}
