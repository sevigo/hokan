package directory

import (
	"bytes"
	"context"
	"encoding/json"
	"path"
	"strings"

	"github.com/nicksnyder/basen"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

var bucketName = "watch:directories"

func New(database core.DB) core.DirectoryStore {
	return &directoryStore{database}
}

type directoryStore struct {
	db core.DB
}

func (s *directoryStore) List(ctx context.Context) ([]*core.Directory, error) {
	var dirs []*core.Directory

	data, err := s.db.ReadBucket(bucketName, &core.ReadBucketOptions{})
	if err != nil {
		if _, ok := err.(*db.ErrBucketNotFound); ok {
			return dirs, nil
		}
		return nil, err
	}

	for _, v := range data {
		dir := &core.Directory{}
		err := json.NewDecoder(strings.NewReader(v)).Decode(dir)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

func (s *directoryStore) Find(ctx context.Context, path int64) (*core.Directory, error) {
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
		return nil, core.ErrDirectoryNotFound
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
	if dir.ID == "" {
		dir.ID = basen.Base62Encoding.EncodeToString([]byte(dir.Path))
	}
	key := dir.ID
	var value bytes.Buffer
	if err := json.NewEncoder(&value).Encode(dir); err != nil {
		return err
	}
	return s.db.Write(bucketName, key, value.String())
}
