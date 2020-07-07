package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/utils/volume"
)

// Storage local
type localStorage struct {
	fileStore         core.FileStore
	bucketName        string
	targetStoragePath string
}

const bufferSize = 1024 * 1024

func New(ctx context.Context, fs core.FileStore, conf core.TargetConfig) (core.TargetStorage, error) {
	if !conf.Active {
		return nil, core.ErrTargetNotActive
	}

	log.WithFields(log.Fields{
		"target": TargetName,
	}).Info("Starting new target storage")

	configurator := NewConfigurator()
	if ok, err := configurator.ValidateSettings(conf.Settings); !ok {
		return nil, err
	}

	return &localStorage{
		targetStoragePath: filepath.Clean(conf.Settings["LOCAL_STORAGE_PATH"]),
		bucketName:        conf.Settings["LOCAL_BUCKET_NAME"],
		fileStore:         fs,
	}, nil
}

func (s *localStorage) Name() string {
	return TargetName
}

func (s *localStorage) Save(ctx context.Context, result chan core.TargetOperationResult, file *core.File, opt *core.TargetStorageSaveOpt) {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})
	logger.Debug("saving file")
	volume := filepath.VolumeName(file.Path)
	base := volume + string(os.PathSeparator)
	relFilePath, err := filepath.Rel(base, file.Path)
	if err != nil {
		result <- core.TargetOperationResultError(err)
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
		result <- core.TargetOperationResultError(errFileSave)
		return
	}
	saveStoreErr := s.fileStore.Save(ctx, TargetName, file)
	result <- core.TargetOperationResultError(saveStoreErr)
	return
}

func (s *localStorage) Restore(ctx context.Context, result chan core.TargetOperationResult, files []*core.File, opt *core.TargetStorageRestoreOpt) {
	log.WithField("target", TargetName).Print("Restore")
	result <- core.TargetOperationResultError(core.ErrNotImplemented)
	return
}

func (s *localStorage) Delete(ctx context.Context, result chan core.TargetOperationResult, files []*core.File, opt *core.TargetStorageDeleteOpt) {
	log.WithField("target", TargetName).Print("Delete")
	result <- core.TargetOperationResultError(core.ErrNotImplemented)
}

// Ping checkes if the local storage is avaible
func (s *localStorage) Ping(ctx context.Context) error {
	if _, err := os.Stat(s.targetStoragePath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *localStorage) Info(ctx context.Context) core.TargetInfo {
	vol := s.targetStoragePath
	if runtime.GOOS == "windows" {
		vol = filepath.VolumeName(s.targetStoragePath)
	}
	f, t := volume.GetVolumeInformation(ctx, vol)
	return core.TargetInfo{
		"free":   fmt.Sprintf("%d", f),
		"total":  fmt.Sprintf("%d", t),
		"volume": vol,
	}
}
