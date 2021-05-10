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

func (f FileInfo) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"name":     f.Name(),
		"size":     f.Size(),
		"mod-time": f.ModTime(),
	})
	// don't care about error here, just return empty json object
	if err != nil {
		return []byte("{}"), nil
	}
	return data, nil
}

type File struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	// if file was renamed
	OldPath  string    `json:"old_path"`
	Checksum string    `json:"checksum"`
	Info     *FileInfo `json:"info"`
	// Targets  []string  `json:"targets"`
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
	Delete(context.Context, string, *File) error
}
