package core

import (
	"context"
	"errors"
)

var ErrDirectoryNotFound = errors.New("directory not found")

type Directory struct {
	ID          string   `json:"id"`
	Active      bool     `json:"active"`
	Path        string   `json:"path"`
	Recursive   bool     `json:"recursive"`
	Machine     string   `json:"machine"`
	IgnoreFiles []string `json:"ignore"`
	Targets     []string `json:"targets"`
}

type DirectoryStats struct {
	Path                string `json:"path"`
	OS                  string `json:"os"`
	TotalFiles          int    `json:"total-files"`
	TotalSubDirectories int    `json:"total-dirs"`
	TotalSize           int64  `json:"total-size"`
}

type DirectoryStore interface {
	List(context.Context) ([]*Directory, error)
	Find(context.Context, int64) (*Directory, error)
	FindName(context.Context, string) (*Directory, error)
	Create(context.Context, *Directory) error
	Update(context.Context, *Directory) error
	Delete(context.Context, *Directory) error
}
