package file

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

var bucketName = "file::snapshot"

type fileStore struct {
	db *db.DB
}

func New(database *db.DB) core.FileStore {
	return &fileStore{database}
}

func (s *fileStore) List(ctx context.Context) ([]*core.File, error) {
	log.Printf("file.List()\n")
	var files []*core.File

	return files, nil
}

func (s *fileStore) Find(ctx context.Context, path string) (*core.File, error) {
	log.Printf("file.List() %q\n", path)
	file := &core.File{}

	return file, nil
}

func (s *fileStore) Save(ctx context.Context, file *core.File) error {
	log.Printf("file.Save() %#v\n", file)

	return nil
}

func (s *fileStore) Delete(ctx context.Context, file *core.File) error {
	log.Printf("file.Delete() %#v\n", file)

	return nil
}
