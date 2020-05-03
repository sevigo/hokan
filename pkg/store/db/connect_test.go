package db_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

func TestConnect(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	dbFile := path.Join(tmpDir, "test.db")
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		want    core.DB
		wantErr bool
	}{
		{
			name:    "case 1",
			path:    dbFile,
			wantErr: false,
		},
		{
			name:    "case 2",
			path:    tmpDir,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Connect(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
