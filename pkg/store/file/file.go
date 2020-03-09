package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

var ErrFileEntryNotFound = errors.New("file entry not found")

type fileStore struct {
	db *db.DB
}

func New(database *db.DB) core.FileStore {
	return &fileStore{database}
}

func (s *fileStore) List(ctx context.Context, bucketName string) ([]*core.File, error) {
	log.Print("file.List()")
	var files []*core.File

	return files, nil
}

func (s *fileStore) Find(ctx context.Context, bucketName, filePath string) (*core.File, error) {
	log.Printf("file.List() %q\n", filePath)
	file := &core.File{}

	key := path.Clean(filePath)
	value, err := s.db.Read(bucketName, key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, ErrFileEntryNotFound
	}

	errJ := json.NewDecoder(bytes.NewReader(value)).Decode(file)
	return file, errJ
}

func (s *fileStore) Save(ctx context.Context, bucketName string, file *core.File) error {
	log.Printf("file.Save() %#v\n", file)
	key := path.Clean(file.Path)
	var value bytes.Buffer
	if err := json.NewEncoder(&value).Encode(file); err != nil {
		return err
	}
	return s.db.Write(bucketName, key, value.String())
}

func (s *fileStore) Delete(ctx context.Context, bucketName string, file *core.File) error {
	log.Printf("file.Delete() %#v\n", file)

	return errors.New("not implemented")
}
