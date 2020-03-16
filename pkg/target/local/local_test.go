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
	filestore "github.com/sevigo/hokan/pkg/store/file"
	"github.com/stretchr/testify/assert"
)

var testFilePath = "local.go"
var bucketName = "test"

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
	fileStore.EXPECT().Find(context.TODO(), TargetName, localPath).Return(nil, filestore.ErrFileEntryNotFound)
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

	tmpFilename := filepath.Join(tmpDir, bucketName, localPath)

	info, err := os.Stat(tmpFilename)
	assert.NoError(t, err)
	assert.Equal(t, testFilePath, info.Name())
}
