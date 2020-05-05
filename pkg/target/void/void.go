package void

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/utils"
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

	return &voidStorage{
		fileStore: fs,
		prefix:    conf.Settings["VOID_PREFIX"],
	}, nil
}

func DefaultConfig() *core.TargetConfig {
	return &core.TargetConfig{
		// always active target for the testing
		Active:      true,
		Name:        TargetName,
		Description: "fake target storage for testing, will print the name of the file",
		Settings: map[string]string{
			"VOID_PREFIX": "",
		},
	}
}

func (s *voidStorage) Save(ctx context.Context, file *core.File, opt *core.TargetStorageSaveOpt) error {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})
	// TODO: this is all the same, move me
	storedFile, err := s.fileStore.Find(ctx, &core.FileSearchOptions{
		FilePath:   file.Path,
		TargetName: TargetName,
	})
	if errors.Is(err, core.ErrFileNotFound) || utils.FileHasChanged(file, storedFile) {
		logger.Debugf("saving the file %s", s.prefix)
		return s.fileStore.Save(ctx, TargetName, file)
	}
	logger.Info("the file has not changed")
	return nil
}

func (s *voidStorage) List(ctx context.Context, opt *core.TargetStorageListOpt) ([]*core.File, error) {
	log.Printf("[void] list\n")
	return nil, nil
}

func (s *voidStorage) Find(ctx context.Context, opt *core.TargetStorageFindOpt) (*core.File, error) {
	log.Printf("[void] find %q\n", opt.Query)
	return nil, nil
}

func (s *voidStorage) Delete(ctx context.Context, file *core.File) error {
	log.Printf("[void] save %#v\n", file)
	return nil
}

func (s *voidStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *voidStorage) Info(ctx context.Context) core.TargetInfo {
	return core.TargetInfo{}
}

func (s *voidStorage) ValidateSettings(settings core.TargetSettings) (bool, error) {
	return true, nil
}
