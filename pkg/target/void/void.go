package void

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/sevigo/hokan/pkg/core"
)

const TargetName = "void"

type voidStorage struct {
	fs core.FileStore
}

func New(fs core.FileStore) (core.TargetStorage, error) {
	return &voidStorage{
		fs: fs,
	}, nil
}

func (s *voidStorage) Save(ctx context.Context, file *core.File) error {
	f, err := s.fs.Find(ctx, TargetName, file.Path)
	if err != nil {
		log.Printf(">>> [void] save a new file [%v] to %q", file, TargetName)
		return s.fs.Save(ctx, TargetName, file)
	}
	if f != nil && f.Checksum == file.Checksum {
		log.Printf("!!! [void] ignore saving, file [%v] already stored in %q", file, TargetName)
		return nil
	}
	if f != nil && f.Checksum != file.Checksum {
		log.Printf("!!! [void] save changed file [%v] to %q", file, TargetName)
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
