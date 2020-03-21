package utils

import "github.com/sevigo/hokan/pkg/core"

func FileHasChanged(newFile, storedFile *core.File) bool {
	if storedFile == nil {
		return true
	}
	if newFile.Checksum != storedFile.Checksum {
		return true
	}
	if newFile.Info != storedFile.Info {
		return true
	}
	return false
}
