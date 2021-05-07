package local

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFile = "testdata/test.txt"

func Test_localStorage_save(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s := &localStorage{
		bucketName:        "test",
		targetStoragePath: tmpDir,
	}

	tests := []struct {
		name     string
		fromPath string
		toPath   string
		wantErr  bool
	}{
		{
			name:     "case 1",
			fromPath: "testdata/test.txt",
			toPath:   path.Join(tmpDir, testFile),
			wantErr:  false,
		},
		{
			name:     "case 2",
			fromPath: "testdata/test.txt",
			toPath:   "",
			wantErr:  true,
		},
		{
			name:     "case 3",
			fromPath: "",
			toPath:   path.Join(tmpDir, testFile),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.save(tt.fromPath, tt.toPath)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			info, err := os.Stat(tt.toPath)
			assert.NoError(t, err)
			assert.Equal(t, int64(11), info.Size())
		})
	}
}
