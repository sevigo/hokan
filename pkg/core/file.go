package core

import "context"

type File struct {
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
