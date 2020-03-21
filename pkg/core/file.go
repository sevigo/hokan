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

type FileListOptions struct {
	TargetName string
	Path       string
	Offset     uint64
	Limit      uint64
}

type FileSearchOptions struct {
	ID         string
	TargetName string
	FilePath   string
}

type FileStore interface {
	List(context.Context, *FileListOptions) ([]*File, error)
	Find(context.Context, *FileSearchOptions) (*File, error)
	Save(context.Context, string, *File) error
	Delete(context.Context, string, *File) error
}
