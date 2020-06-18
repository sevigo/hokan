package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/utils/volume"
)

var bucketNameRegexp = regexp.MustCompile("^[a-zA-Z0-9_.]+$")

const TargetName = "local"

var config = core.TargetConfig{
	Active:      false,
	Name:        TargetName,
	Description: "store the files on the local disk",
	Settings: core.TargetSettings{
		"LOCAL_STORAGE_PATH": "",
		"LOCAL_BUCKET_NAME":  "",
	},
}

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

	s := &localStorage{}
	if ok, err := s.ValidateSettings(conf.Settings); !ok {
		return nil, err
	}

	s.targetStoragePath = filepath.Clean(conf.Settings["LOCAL_STORAGE_PATH"])
	s.bucketName = conf.Settings["LOCAL_BUCKET_NAME"]
	s.fileStore = fs

	return s, nil
}

func (s *localStorage) DefaultConfig() *core.TargetConfig {
	return &config
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

func (s *localStorage) ValidateSettings(settings core.TargetSettings) (bool, error) {
	logger := log.WithField("target", TargetName)
	logger.Infof("ValidateSettings(): %+v", settings)

	path, ok := settings["LOCAL_STORAGE_PATH"]
	if !ok {
		return false, fmt.Errorf("LOCAL_STORAGE_PATH is empty")
	}
	if _, err := os.Stat(filepath.Clean(path)); os.IsNotExist(err) {
		return false, fmt.Errorf("%q does not exist", filepath.Clean(path))
	}

	bucket, ok := settings["LOCAL_BUCKET_NAME"]
	if !ok {
		return false, fmt.Errorf("LOCAL_BUCKET_NAME is empty")
	}
	if bucket == "" {
		return false, fmt.Errorf("LOCAL_BUCKET_NAME is empty")
	}

	match := bucketNameRegexp.MatchString(bucket)
	if !match {
		return false, fmt.Errorf("bucket name contains illegal characters")
	}

	return true, nil
}
