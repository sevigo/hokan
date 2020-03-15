package core

import "context"

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

type TargetRegister interface {
	AllTargets() map[string]TargetFactory
	GetTarget(name string) TargetStorage
}

// TargetFactory is a function that returns a TargetStorage.
type TargetFactory func(context.Context, FileStore) (TargetStorage, error)
