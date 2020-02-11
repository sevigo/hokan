package main

import (
	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
	"github.com/sevigo/hokan/pkg/store/directory"
)

var storeSet = wire.NewSet(
	provideDatabase,
	provideDirectoryStore,
)

func provideDatabase(config config.Config) (*db.DB, error) {
	return db.Connect(config.Database.Path)
}

func provideDirectoryStore(db *db.DB) core.DirectoryStore {
	dirs := directory.New(db)
	return dirs
}
