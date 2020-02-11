package directory

import (
	"context"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

func New(db *db.DB) core.DirectoryStore {
	return &directoryStore{db}
}

type directoryStore struct {
	db *db.DB
}

func (s *directoryStore) List(ctx context.Context, id int64) ([]*core.Directory, error) {
	var dirs []*core.Directory

	return dirs, nil
}

func (s *directoryStore) Find(ctx context.Context, id int64) (*core.Directory, error) {
	dir := &core.Directory{}

	return dir, nil
}

func (s *directoryStore) FindName(ctx context.Context, id int64, name string) (*core.Directory, error) {
	dir := &core.Directory{}

	return dir, nil
}

func (s *directoryStore) Update(ctx context.Context, dir *core.Directory) error {
	return nil
}

func (s *directoryStore) Delete(ctx context.Context, dir *core.Directory) error {
	return nil
}

func (s *directoryStore) Create(ctx context.Context, dir *core.Directory) error {
	return nil
}
