package core

import "context"

type Directory struct {
	Active      bool
	Path        string
	Recursive   bool
	Machine     string
	IgnoreFiles []string
	// 	TODO: rename to Targets []string
	Target []string
}

type DirectoryStore interface {
	List(context.Context) ([]*Directory, error)

	Find(context.Context, int64) (*Directory, error)

	FindName(context.Context, string) (*Directory, error)

	Create(context.Context, *Directory) error

	Update(context.Context, *Directory) error

	Delete(context.Context, *Directory) error
}
