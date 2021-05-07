package void

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

const name = "void"

type voidBackup struct {
	fileStore core.FileStore
	prefix    string
}

func New(_ context.Context, fs core.FileStore, options *core.BackupOptions) (core.Backup, error) {
	log.WithFields(log.Fields{
		"backup": name,
	}).Info("Starting new connection to the backup storage")

	return &voidBackup{
		fileStore: fs,
		prefix:    "void",
	}, nil
}

func (v *voidBackup) Name() string {
	return name
}

func (v *voidBackup) Save(ctx context.Context, result chan core.BackupResult, file *core.File, _ *core.BackupOperationOptions) {
	logger := log.WithFields(log.Fields{
		"backup": name,
		"file":   file.Path,
	})
	logger.Infof("void.Save(): saving the file %q", file.Path)
	err := v.fileStore.Save(ctx, name, file)
	if err != nil {
		result <- core.BackupResult{
			Success: false,
			Error:   err,
		}
	} else {
		result <- core.BackupResult{
			Success: true,
			Message: fmt.Sprintf("%s for file %q", core.BackupSuccessMessage, file.Path),
		}
	}
}

func (v *voidBackup) Restore(ctx context.Context, result chan core.BackupResult, files []*core.File, _ *core.BackupOperationOptions) {
	log.Printf("[void] restore %#v\n", files)
	result <- core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}
}

func (v *voidBackup) Delete(ctx context.Context, result chan core.BackupResult, files []*core.File, _ *core.BackupOperationOptions) {
	log.Printf("[void] delete %#v\n", files)
	result <- core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}
}

func (v *voidBackup) Ping(ctx context.Context) error {
	return nil
}
