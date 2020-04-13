package core

import (
	"context"
	"errors"
)

type TargetStorageStatus int

var ErrTargetNotActive = errors.New("target is not active")

var ErrTargetConfigNotFound = errors.New("default config for target not found")

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
	AllConfigs() map[string]*TargetConfig
	AllTargets() map[string]TargetFactory
	GetTarget(name string) TargetStorage
	GetConfig(context.Context, string) (*TargetConfig, error)
	SetConfig(context.Context, *TargetConfig) error
}

type TargetConfig struct {
	Active      bool
	Name        string
	Description string
	Settings    map[string]string
	Stats       TargetStats
}

type TargetStats struct {
	TotalFiles uint64
	TotalSize  uint64
}

// TargetFactory is a function that returns a TargetStorage.
type TargetFactory func(context.Context, FileStore, TargetConfig) (TargetStorage, error)
