package core

import (
	"context"
	"errors"
)

var ErrDirectoryNotFound = errors.New("directory not found")

type Directory struct {
	ID          string
	Active      bool
	Path        string
	Recursive   bool
	Machine     string
	IgnoreFiles []string
	Targets     []string
}

type DirectoryStore interface {
	List(context.Context) ([]*Directory, error)
	Find(context.Context, int64) (*Directory, error)
	FindName(context.Context, string) (*Directory, error)
	Create(context.Context, *Directory) error
	Update(context.Context, *Directory) error
	Delete(context.Context, *Directory) error
}
