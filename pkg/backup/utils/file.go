package utils

import (
	"github.com/sevigo/hokan/pkg/core"
)

func FileHasChanged(newFile, storedFile *core.File) bool {
	if storedFile == nil {
		return true
	}
	if newFile.Path != storedFile.Path {
		return true
	}
	if newFile.Checksum != storedFile.Checksum {
		return true
	}
	return false
}

func sizeHasChanged(newFile, storedFile *core.File) bool {
	if newFile.Info != nil && storedFile.Info != nil && newFile.Info.Size() != storedFile.Info.Size() {
		return true
	}
	return false
}
