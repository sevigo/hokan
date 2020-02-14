package directory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"path"
	"strings"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

var bucketName = "watch:directories"

func New(db *db.DB) core.DirectoryStore {
	return &directoryStore{db}
}

type directoryStore struct {
	db *db.DB
}

func (s *directoryStore) List(ctx context.Context) ([]*core.Directory, error) {
	var dirs []*core.Directory

	data, err := s.db.ReadBucket(bucketName)
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		dir := &core.Directory{}
		json.NewDecoder(strings.NewReader(v)).Decode(dir)
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

func (s *directoryStore) Find(ctx context.Context, id int64) (*core.Directory, error) {
	dir := &core.Directory{}

	return dir, nil
}

func (s *directoryStore) FindName(ctx context.Context, name string) (*core.Directory, error) {
	log.Printf("directory.FindName(): %s\n", name)
	dir := &core.Directory{}

	key := path.Clean(name)
	value, err := s.db.Read(bucketName, key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, errors.New("Entry was not found")
	}

	err = json.NewDecoder(bytes.NewReader(value)).Decode(dir)
	return dir, err
}

func (s *directoryStore) Update(ctx context.Context, dir *core.Directory) error {
	return nil
}

func (s *directoryStore) Delete(ctx context.Context, dir *core.Directory) error {
	return nil
}

func (s *directoryStore) Create(ctx context.Context, dir *core.Directory) error {
	log.Printf("directory.Create(): %#v\n", dir)
	key := path.Clean(dir.Path)
	var value bytes.Buffer
	if err := json.NewEncoder(&value).Encode(dir); err != nil {
		return err
	}
	return s.db.Write(bucketName, key, value.String())
}
