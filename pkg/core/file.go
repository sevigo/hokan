package core

import (
	"context"
	"errors"
)

var ErrFileNotFound = errors.New("file not found")

type File struct {
	ID       string
	Path     string
	Checksum string
	Info     string
	Targets  []string
}

type FileStore interface {
	List(context.Context, string) ([]*File, error)
	Find(context.Context, string, string) (*File, error)
	Save(context.Context, string, *File) error
	Delete(context.Context, string, *File) error
}
