package backup

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/backup/local"
	"github.com/sevigo/hokan/pkg/backup/minio"
	"github.com/sevigo/hokan/pkg/backup/void"
	"github.com/sevigo/hokan/pkg/core"
)

var backupStorageRegister = map[string]core.BackupFactory{
	"void":  void.New,
	"local": local.New,
	"minio": minio.New,
}

func New(ctx context.Context, fileStore core.FileStore, options *core.BackupOptions) (core.Backup, error) {
	backup, err := getBackupStorage(options.Name)
	if err != nil {
		log.WithField("backup", options.Name).
			WithError(err).
			Fatal("can't initiate new backup storage")
		return nil, err
	}
	return backup(ctx, fileStore, options)
}

func getBackupStorage(name string) (core.BackupFactory, error) {
	if _, ok := backupStorageRegister[name]; ok {
		return backupStorageRegister[name], nil
	}
	return nil, core.ErrBackupStorageNotFound
}
