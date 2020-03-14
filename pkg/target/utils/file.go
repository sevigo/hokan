package utils

import "github.com/sevigo/hokan/pkg/core"

func FileHasChanged(newFile, storedFile *core.File) bool {
	if storedFile == nil {
		return true
	}
	if newFile.Checksum != storedFile.Checksum {
		return true
	}
	return false
}
