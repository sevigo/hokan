package local

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

var testFilePath = "testdata/test.txt"
var bucketName = "test"

func TestLocalStorageSaveNew(t *testing.T) {
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

func TestLocalStorageSaveWrongPath(t *testing.T) {
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
	assert.Contains(t, data.Message, "cannot open file [/wrong/path]")
}

func TestLocalStorageSaveError(t *testing.T) {
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
