package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/nicksnyder/basen"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

type fileStore struct {
	db *db.DB
}

func New(database *db.DB) core.FileStore {
	return &fileStore{database}
}

func (s *fileStore) List(ctx context.Context, opt *core.FileListOptions) ([]*core.File, error) {
	log.Print("file.List()")
	if opt == nil {
		return nil, fmt.Errorf("empty list options")
	}
	var files []*core.File

	data, err := s.db.ReadBucket(opt.TargetName, &db.ReadBucketOptions{})
	if errors.Is(err, &db.ErrBucketNotFound{}) {
		return nil, core.ErrTargetNotActive
	}
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		file := &core.File{}
		err := json.NewDecoder(strings.NewReader(v)).Decode(file)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (s *fileStore) Find(ctx context.Context, opt *core.FileSearchOptions) (*core.File, error) {
	if opt == nil {
		return nil, fmt.Errorf("empty search options")
	}
	file := &core.File{}

	var key string
	if opt.FilePath != "" {
		key = path.Clean(opt.FilePath)
	} else if opt.ID != "" {
		filePath, err := basen.Base62Encoding.DecodeString(opt.ID)
		if err != nil {
			return nil, err
		}
		key = path.Clean(string(filePath))
	} else {
		return nil, fmt.Errorf("empty search parameter")
	}

	value, err := s.db.Read(opt.TargetName, key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, core.ErrFileNotFound
	}

	errJ := json.NewDecoder(bytes.NewReader(value)).Decode(file)
	return file, errJ
}

func (s *fileStore) Save(ctx context.Context, bucketName string, file *core.File) error {
	log.Printf("file.Save() %q\n", file.Path)
	key := path.Clean(file.Path)
	if file.ID == "" {
		file.ID = basen.Base62Encoding.EncodeToString([]byte(key))
	}
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
