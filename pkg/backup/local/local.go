package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

const name = "local"

// Storage local
type localStorage struct {
	fileStore         core.FileStore
	bucketName        string
	targetStoragePath string
}

const bufferSize = 1024 * 1024

func New(ctx context.Context, fs core.FileStore, conf *core.BackupOptions) (core.Backup, error) {
	log.WithFields(log.Fields{
		"backup": name,
	}).Info("Starting new backup storage")

	return &localStorage{
		targetStoragePath: filepath.Clean(conf.TargetURL),
		bucketName:        name,
		fileStore:         fs,
	}, nil
}

func (s *localStorage) Name() string {
	return name
}

func (s *localStorage) Save(ctx context.Context, result chan core.BackupResult, file *core.File, opt *core.BackupOperationOptions) {
	logger := log.WithFields(log.Fields{
		"backup": name,
		"file":   file.Path,
	})
	logger.Debug("saving file")
	volume := filepath.VolumeName(file.Path)
	base := volume + string(os.PathSeparator)
	relFilePath, err := filepath.Rel(base, file.Path)
	if err != nil {
		result <- core.BackupResult{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("can't parse storage path: %q, base: %q", file.Path, base),
		}
		return
	}
	// on windows volume will be 'C:', so we just remove :
	// on all other systems it will be empty
	if volume != "" {
		volume = strings.TrimSuffix(volume, ":")
	}
	to := filepath.Join(s.targetStoragePath, s.bucketName, volume, relFilePath)
	errFileSave := s.save(file.Path, to)
	if errFileSave != nil {
		result <- core.BackupResult{
			Success: false,
			Error:   errFileSave,
			Message: fmt.Sprintf("can't save file %q to %q", file.Path, to),
		}
		return
	}

	saveStoreErr := s.fileStore.Save(ctx, name, file)
	if saveStoreErr != nil {
		result <- core.BackupResult{
			Success: false,
			Error:   saveStoreErr,
			Message: fmt.Sprintf("can't save backup info for the file %q to %q", file.Path, to),
		}
		return
	}

	// backup was successful
	result <- core.BackupResult{
		Success: true,
		Message: core.BackupSuccessMessage,
	}
}

func (s *localStorage) Restore(ctx context.Context, result chan core.BackupResult, files []*core.File, opt *core.BackupOperationOptions) {
	log.WithField("backup", name).Print("calling restore()")
	result <- core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}
}

func (s *localStorage) Delete(ctx context.Context, result chan core.BackupResult, files []*core.File, opt *core.BackupOperationOptions) {
	log.WithField("backup", name).Print("calling delete()")
	result <- core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}
}

// Ping checkes if the local storage is avaible
func (s *localStorage) Ping(ctx context.Context) error {
	if _, err := os.Stat(s.targetStoragePath); os.IsNotExist(err) {
		return err
	}
	return nil
}
