package main

import (
	"context"

	"github.com/google/wire"

	"github.com/sevigo/hokan/cmd/hokan/config"
	"github.com/sevigo/hokan/pkg/backup"
	"github.com/sevigo/hokan/pkg/core"
)

var backupSet = wire.NewSet(
	provideBackup,
)

func provideBackup(ctx context.Context, fileStore core.FileStore, events core.EventCreator, config config.Config) (core.Backup, error) {
	return backup.New(ctx, fileStore, events, &core.BackupOptions{
		Name:      config.Backup.Name,
		TargetURL: config.Backup.TargetLocalPath,
	})
}
