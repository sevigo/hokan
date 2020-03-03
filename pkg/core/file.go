package core

import "context"

type File struct {
	Path    string
	Targets []string
}

type FileStore interface {
	List(context.Context) ([]*File, error)

	Find(context.Context, string) (*File, error)

	Save(context.Context, *File) error

	Delete(context.Context, *File) error
}
