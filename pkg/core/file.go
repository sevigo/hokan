package core

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

var ErrFileNotFound = errors.New("file not found")

type FileInfo struct {
	os.FileInfo
}

func (f FileInfo) JSON() string {
	data, err := json.Marshal(map[string]interface{}{
		"Name":    f.Name(),
		"Size":    f.Size(),
		"Mode":    f.Mode(),
		"ModTime": f.ModTime(),
	})
	if err != nil {
		return ""
	}
	return string(data)
}

type File struct {
	ID       string
	Path     string
	Checksum string
	Info     *FileInfo
	Targets  []string
}

type FileListOptions struct {
	TargetName string
	Path       string
	Offset     uint64
	Limit      uint64
}

type FileSearchOptions struct {
	ID         string
	TargetName string
	FilePath   string
}

type FileStore interface {
	List(context.Context, *FileListOptions) ([]*File, error)
	Find(context.Context, *FileSearchOptions) (*File, error)
	Save(context.Context, string, *File) error
	Update(context.Context, string, *File) error
	Delete(context.Context, string, *File) error
}
