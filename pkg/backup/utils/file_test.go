package utils

import (
	"testing"

	"github.com/sevigo/hokan/pkg/core"
)

func TestFileHasChanged(t *testing.T) {
	tests := []struct {
		name       string
		newFile    *core.File
		storedFile *core.File
		want       bool
	}{
		{
			name: "case 1",
			newFile: &core.File{
				Path:     "/test/file.txt",
				Checksum: "abc",
			},
			storedFile: nil,
			want:       true,
		},
		{
			name: "case 2",
			newFile: &core.File{
				Path:     "/test/file.txt",
				Checksum: "abc",
			},
			storedFile: &core.File{
				Path:     "/test/file.txt",
				Checksum: "abX",
			},
			want: true,
		},
		{
			name: "case 2",
			newFile: &core.File{
				Path:     "/test/file.txt",
				Checksum: "abc",
			},
			storedFile: &core.File{
				Path:     "/test/file.txt",
				Checksum: "abc",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileHasChanged(tt.newFile, tt.storedFile); got != tt.want {
				t.Errorf("FileHasChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}
