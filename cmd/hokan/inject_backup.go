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
	proviceBackupResultChan,
)

func provideBackup(ctx context.Context,
	fileStore core.FileStore,
	events core.EventCreator,
	results chan core.BackupResult,
	config config.Config) (core.Backup, error) {
	return backup.New(ctx, fileStore, events, results, &core.BackupOptions{
		Name:      config.Backup.Name,
		TargetURL: config.Backup.TargetLocalPath,
	})
}

func proviceBackupResultChan() chan core.BackupResult {
	return make(chan core.BackupResult)
}
