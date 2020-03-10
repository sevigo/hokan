package void

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

const TargetName = "void"

type voidStorage struct {
	fs core.FileStore
}

func New(ctx context.Context, fs core.FileStore) (core.TargetStorage, error) {
	return &voidStorage{
		fs: fs,
	}, nil
}

func (s *voidStorage) Save(ctx context.Context, file *core.File) error {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	f, err := s.fs.Find(ctx, TargetName, file.Path)
	if err != nil {
		logger.Debug("saving new file")
		return s.fs.Save(ctx, TargetName, file)
	}
	if f != nil && f.Checksum == file.Checksum {
		logger.Debug("ignore, file already stored")
		return nil
	}
	if f != nil && f.Checksum != file.Checksum {
		logger.Debug("save changed file")
		return s.fs.Save(ctx, TargetName, file)
	}

	return nil
}

func (s *voidStorage) List(context.Context) ([]*core.File, error) {
	log.Printf("[void] list\n")
	return nil, nil
}

func (s *voidStorage) Find(ctx context.Context, q string) (*core.File, error) {
	log.Printf("[void] find %q\n", q)
	return nil, nil
}

func (s *voidStorage) Delete(ctx context.Context, file *core.File) error {
	log.Printf("[void] save %#v\n", file)
	return nil
}

func (s *voidStorage) Ping(ctx context.Context) error {
	return nil
}
