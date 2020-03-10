package core

import "context"

type Target interface {
	// InitTargets(ctx)
}

type TargetStorageStatus int

const (
	TargetStorageOK TargetStorageStatus = iota
	TargetStoragePaused
	TargetStorageError
)

type TargetStorage interface {
	List(context.Context) ([]*File, error)
	Find(context.Context, string) (*File, error)
	Save(context.Context, *File) error
	Delete(context.Context, *File) error
	Ping(context.Context) error
}

// TargetFactory is a function that returns a TargetStorage.
type TargetFactory func(context.Context, FileStore) (TargetStorage, error)
