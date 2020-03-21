package void

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

var testFilePath = "/test/file.txt"

func TestConfig(t *testing.T) {
	conf := DefaultConfig()
	assert.Equal(t, "void", conf.Name)
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

	file := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   testFilePath,
	}).Return(nil, core.ErrFileNotFound)
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

	store := &voidStorage{
		fileStore: fileStore,
	}

	err := store.Save(context.TODO(), file)
	assert.NoError(t, err)
}

func Test_voidStorageSaveChanged(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileA := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
	}
	fileB := &core.File{
		Path:     testFilePath,
		Checksum: "abX",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   testFilePath,
	}).Return(fileA, nil)
	fileStore.EXPECT().Save(context.TODO(), TargetName, fileB).Return(nil)

	store := &voidStorage{
		fileStore: fileStore,
	}

	err := store.Save(context.TODO(), fileB)
	assert.NoError(t, err)
}

func Test_minioStore_NoSave(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileA := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   testFilePath,
	}).Return(fileA, nil)

	store := &voidStorage{
		fileStore: fileStore,
	}

	err := store.Save(context.TODO(), fileA)
	assert.NoError(t, err)
}
