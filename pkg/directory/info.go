package directory

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/sevigo/hokan/pkg/core"
)

func Stats(path string) (*core.DirectoryStats, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, core.ErrDirectoryNotFound
	}
	stats := &core.DirectoryStats{}
	stats.OS = runtime.GOOS
	stats.Path = path
	err := filepath.Walk(path, func(absoluteFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			stats.TotalSubDirectories++
		} else {
			stats.TotalFiles++
			stats.TotalSize += fileInfo.Size()
		}
		return nil
	})

	return stats, err
}
