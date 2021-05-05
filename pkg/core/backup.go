package core

import (
	"context"
	"errors"
)

var ErrBackupStorageNotFound = errors.New("Backup storage not found")

const BackupSuccessMessage = "requested operation was successful"

type BackupOptions struct {
	Name            string
	TargetURL       string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type BackupOperationOptions struct {
}

type BackupResult struct {
	Success bool
	Error   error
	Message string
}

// BackupFactory is a function that returns a Backup.
type BackupFactory func(ctx context.Context, fs FileStore, options *BackupOptions) (Backup, error)

type Backup interface {
	Name() string
	Save(ctx context.Context, result chan BackupResult, file *File, opt *BackupOperationOptions)
	Restore(ctx context.Context, result chan BackupResult, files []*File, opt *BackupOperationOptions)
	Delete(ctx context.Context, result chan BackupResult, files []*File, opt *BackupOperationOptions)
	Ping(context.Context) error
}
