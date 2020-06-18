package local

import (
	"context"
	"errors"
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
	store := &localStorage{}
	conf := store.DefaultConfig()
	assert.Equal(t, "local", conf.Name)
	assert.Equal(t, false, conf.Active)
}

func TestNewNotActive(t *testing.T) {
	store := &localStorage{}
	conf := store.DefaultConfig()
	_, err := New(context.Background(), nil, *conf)
	assert.EqualError(t, err, "target is not active")
}

func TestNewActive(t *testing.T) {
	store := &localStorage{}
	conf := store.DefaultConfig()
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
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := &localStorage{
		bucketName:        bucketName,
		targetStoragePath: tmpDir,
		fileStore:         fileStore,
	}

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, file, &core.TargetStorageSaveOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: true,
		Error:   nil,
		Message: "requested operation was successful",
	}, data)
}

func TestVoidStorageSaveWrongPath(t *testing.T) {
	file := &core.File{
		Path:     "/wrong/path",
		Checksum: "abc",
	}
	store := &localStorage{
		bucketName: bucketName,
	}

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, file, &core.TargetStorageSaveOpt{})
	data := <-result
	assert.Equal(t, false, data.Success)
	assert.Equal(t, "cannot open file  [/wrong/path]: open /wrong/path: The system cannot find the path specified.", data.Message)
}

func TestVoidStorageSaveError(t *testing.T) {
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
	errTest := errors.New("test")
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(errTest)

	store := &localStorage{
		bucketName: bucketName,
		fileStore:  fileStore,
	}

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, file, &core.TargetStorageSaveOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   errTest,
		Message: "test",
	}, data)
}

func TestRestore(t *testing.T) {
	store := &localStorage{}
	result := make(chan core.TargetOperationResult)
	go store.Restore(context.TODO(), result, []*core.File{}, &core.TargetStorageRestoreOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   core.ErrNotImplemented,
		Message: "not implemented",
	}, data)
}

func TestDelete(t *testing.T) {
	store := &localStorage{}
	result := make(chan core.TargetOperationResult)
	go store.Delete(context.TODO(), result, []*core.File{}, &core.TargetStorageDeleteOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   core.ErrNotImplemented,
		Message: "not implemented",
	}, data)
}

func TestVoidStoragePing(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	localPath := filepath.Join(pwd, testFilePath)
	fileStore := mocks.NewMockFileStore(controller)

	store := &localStorage{
		bucketName:        bucketName,
		fileStore:         fileStore,
		targetStoragePath: localPath,
	}

	err = store.Ping(context.TODO())
	assert.NoError(t, err)
}

func TestVoidStoragePingError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := &localStorage{
		bucketName:        bucketName,
		targetStoragePath: "/wrong/path",
	}

	err := store.Ping(context.TODO())
	assert.Error(t, err)
}

func TestLocalInfo(t *testing.T) {
	store := &localStorage{}
	conf := store.DefaultConfig()
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
	store := &localStorage{}
	conf := store.DefaultConfig()
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
