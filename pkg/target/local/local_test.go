package local

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
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
