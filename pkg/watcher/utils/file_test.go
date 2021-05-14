package utils

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validPath = path.Join("testdata", "test.txt")

const errorPath = "testdata/nofile"

func TestFileChecksumInfo(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantSum  string
		wantSize int64
		wantErr  bool
	}{
		{
			name:     "case 1",
			path:     validPath,
			wantErr:  false,
			wantSize: int64(11),
			wantSum:  "ea1f27bdefa157182ae6f08e468f91c826525cdea7a9f42783330b2d63e89958",
		},
		{
			name:    "case 2",
			path:    errorPath,
			wantErr: true,
		},
		{
			name:    "case 3",
			path:    "testdata",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, info, err := FileChecksumInfo(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantSum, sum)
			assert.Equal(t, tt.wantSize, info.Size())
		})
	}
}
