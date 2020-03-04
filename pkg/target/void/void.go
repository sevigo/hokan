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
	log.Printf("New void storage created\n")
	return &voidStorage{
		fs: fs,
	}, nil
}

func (s *voidStorage) Save(ctx context.Context, file *core.File) error {
	log.Printf("[void] save %#v\n", file)
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
