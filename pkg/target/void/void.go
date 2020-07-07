package void

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

const TargetName = "void"

type voidStorage struct {
	fileStore core.FileStore
	prefix    string
}

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

	return &voidStorage{
		fileStore: fs,
		prefix:    conf.Settings["VOID_PREFIX"],
	}, nil
}

func (s *voidStorage) Name() string {
	return TargetName
}

func (s *voidStorage) Save(ctx context.Context, result chan core.TargetOperationResult, file *core.File, opt *core.TargetStorageSaveOpt) {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})
	logger.Debugf("saving the file %s", s.prefix)
	saveErr := s.fileStore.Save(ctx, TargetName, file)
	result <- core.TargetOperationResultError(saveErr)
}

func (s *voidStorage) Restore(ctx context.Context, result chan core.TargetOperationResult, files []*core.File, opt *core.TargetStorageRestoreOpt) {
	result <- core.TargetOperationResultError(core.ErrNotImplemented)
}

func (s *voidStorage) Delete(ctx context.Context, result chan core.TargetOperationResult, files []*core.File, opt *core.TargetStorageDeleteOpt) {
	log.Printf("[void] save %#v\n", files)
	result <- core.TargetOperationResultError(core.ErrNotImplemented)
}

func (s *voidStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *voidStorage) Info(ctx context.Context) core.TargetInfo {
	return core.TargetInfo{}
}
