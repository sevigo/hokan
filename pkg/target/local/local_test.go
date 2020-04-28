package local

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

var testFilePath = "local.go"
var bucketName = "test"

func TestConfig(t *testing.T) {
	conf := DefaultConfig()
	assert.Equal(t, "local", conf.Name)
	assert.Equal(t, false, conf.Active)
}

func TestNewNotActive(t *testing.T) {
	conf := DefaultConfig()
	_, err := New(context.Background(), nil, *conf)
	assert.EqualError(t, err, "target is not active")
}

func TestNewActive(t *testing.T) {
	conf := DefaultConfig()
	conf.Active = true
	conf.Settings["LOCAL_BUCKET_NAME"] = "test"
	conf.Settings["LOCAL_STORAGE_PATH"] = "."
	_, err := New(context.Background(), nil, *conf)
	assert.NoError(t, err)
}

func Test_voidStorageSaveNew(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	localPath := filepath.Join(pwd, testFilePath)

	file := &core.File{
		Path:     localPath,
		Checksum: "abc",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   localPath,
	}).Return(nil, core.ErrFileNotFound)
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := &localStorage{
		bucketName:        bucketName,
		targetStoragePath: tmpDir,
		fileStore:         fileStore,
	}

	err = store.Save(context.TODO(), file)
	assert.NoError(t, err)
}

func Test_voidStorageSaveNoChanges(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	localPath := filepath.Join(pwd, testFilePath)

	file := &core.File{
		Path:     localPath,
		Checksum: "abc",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   localPath,
	}).Return(file, nil)

	store := &localStorage{
		bucketName: bucketName,
		fileStore:  fileStore,
	}

	err = store.Save(context.TODO(), file)
	assert.NoError(t, err)
}

func TestLocalInfo(t *testing.T) {
	conf := DefaultConfig()
	conf.Active = true

	if runtime.GOOS == "windows" {
		conf.Settings["LOCAL_STORAGE_PATH"] = "C:\\"
		conf.Settings["LOCAL_BUCKET_NAME"] = "test"
	} else {
		pwd, err := os.Getwd()
		assert.NoError(t, err)
		conf.Settings["LOCAL_STORAGE_PATH"] = pwd
		conf.Settings["LOCAL_BUCKET_NAME"] = "test"
	}

	target, err := New(context.Background(), nil, *conf)
	assert.NoError(t, err)
	info := target.Info(context.TODO())
	assert.NotEmpty(t, info)
	assert.NotEmpty(t, info["total"])
	assert.NotEmpty(t, info["free"])
	assert.NotEmpty(t, info["volume"])
}
func TestNewError(t *testing.T) {
	conf := DefaultConfig()
	conf.Active = true
	conf.Settings["LOCAL_STORAGE_PATH"] = ""
	conf.Settings["LOCAL_BUCKET_NAME"] = ""

	_, err := New(context.Background(), nil, *conf)
	assert.Error(t, err)
}

func Test_localStorage_ValidateSettings(t *testing.T) {
	store := &localStorage{}
	pwd, err := os.Getwd()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		settings core.TargetSettings
		wantErr  bool
	}{
		{
			name: "case 1",
			settings: core.TargetSettings{
				"LOCAL_STORAGE_PATH": pwd,
				"LOCAL_BUCKET_NAME":  "test",
			},
			wantErr: false,
		},
		{
			name: "case 2",
			settings: core.TargetSettings{
				"LOCAL_STORAGE_PATH": pwd,
			},
			wantErr: true,
		},
		{
			name: "case 3",
			settings: core.TargetSettings{
				"LOCAL_BUCKET_NAME": "test",
			},
			wantErr: true,
		},
		{
			name: "case 4",
			settings: core.TargetSettings{
				"LOCAL_STORAGE_PATH": "/not_valid_path",
				"LOCAL_BUCKET_NAME":  "test.me",
			},
			wantErr: true,
		},
		{
			name: "case 5",
			settings: core.TargetSettings{
				"LOCAL_STORAGE_PATH": pwd,
				"LOCAL_BUCKET_NAME":  ");DROP",
			},
			wantErr: true,
		},
		{
			name: "case 6",
			settings: core.TargetSettings{
				"LOCAL_STORAGE_PATH": pwd,
				"LOCAL_BUCKET_NAME":  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.ValidateSettings(tt.settings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
