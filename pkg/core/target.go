package core

import "context"

type Target struct{}

type TargetStorage interface {
	List(context.Context) ([]*File, error)

	Find(context.Context, string) (*File, error)

	Save(context.Context, *File) error

	Delete(context.Context, *File) error
}

// TargetFactory is a function that returns a TargetStorage.
type TargetFactory func(*FileStore) (TargetStorage, error)
