package core

import (
	"context"
	"errors"
)

type TargetStorageStatus int

var ErrTargetNotActive = errors.New("target is not active")

var ErrTargetConfigNotFound = errors.New("default config for target not found")

const defaulSuccesstMessage = "requested operation was successful"

const (
	TargetStorageOK TargetStorageStatus = iota
	TargetStoragePaused
	TargetStorageError
)

type TargetStorageDeleteOpt struct {
}

type TargetStorageSaveOpt struct {
	Force bool
}

type TargetStorageFindOpt struct {
	Query string
}

type TargetStorageListOpt struct {
}

type TargetStorageRestoreOpt struct {
	LocalPath         string
	OverrideOriginals bool
	UseOriginalPath   bool
}

type TargetOperationResult struct {
	Success bool
	Error   error
	Message string
}

func TargetOperationResultError(err error) TargetOperationResult {
	return targetOperationResultChan(err, "")
}

func TargetOperationResultSuccess(msg string) TargetOperationResult {
	return targetOperationResultChan(nil, msg)
}

func targetOperationResultChan(err error, msg string) TargetOperationResult {
	result := TargetOperationResult{}
	if err != nil {
		if msg == "" {
			msg = err.Error()
		}
		result = TargetOperationResult{
			Success: false,
			Message: msg,
			Error:   err,
		}
	} else {
		if msg == "" {
			msg = defaulSuccesstMessage
		}
		result = TargetOperationResult{
			Success: true,
			Message: msg,
		}
	}
	return result
}

type TargetStorage interface {
	DefaultConfig() *TargetConfig
	Save(ctx context.Context, result chan TargetOperationResult, file *File, opt *TargetStorageSaveOpt)
	Restore(ctx context.Context, result chan TargetOperationResult, files []*File, opt *TargetStorageRestoreOpt)
	Delete(ctx context.Context, result chan TargetOperationResult, files []*File, opt *TargetStorageDeleteOpt)
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
