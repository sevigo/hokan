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

type TargetStorageSaveOpt struct {
	Force bool
}

type TargetStorageFindOpt struct {
	Query string
}

type TargetStorageListOpt struct {
}

type TargetStorage interface {
	List(ctx context.Context, opt *TargetStorageListOpt) ([]*File, error)
	Find(ctx context.Context, opt *TargetStorageFindOpt) (*File, error)
	Save(ctx context.Context, file *File, opt *TargetStorageSaveOpt) error
	Delete(context.Context, *File) error
	Ping(context.Context) error
	Info(context.Context) TargetInfo
	ValidateSettings(settings TargetSettings) (bool, error)
}

type TargetRegister interface {
	AllConfigs() map[string]TargetConfig
	AllTargets() map[string]Target
	GetTarget(name string) TargetStorage
	GetConfig(context.Context, string) (*TargetConfig, error)
	SetConfig(context.Context, *TargetConfig) error
}

type TargetConfig struct {
	Active      bool           `json:"active"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Settings    TargetSettings `json:"settings"`
}

type Target struct {
	Name   string              `json:"name"`
	Status TargetStorageStatus `json:"status"`
	Info   TargetInfo          `json:"info"`
}

type TargetInfo map[string]string

type TargetSettings map[string]string

// TargetFactory is a function that returns a TargetStorage.
type TargetFactory func(context.Context, FileStore, TargetConfig) (TargetStorage, error)
