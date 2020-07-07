package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const validPath = "testdata/test.txt"
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
			wantSum:  "5e2bf57d3f40c4b6df69daf1936cb766f832374b4fc0259a7cbff06e2f70f269",
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
