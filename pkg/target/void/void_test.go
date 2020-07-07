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

func Test_voidStorageSaveNew(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	file := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

	store := &voidStorage{
		fileStore: fileStore,
	}

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, file, nil)
	msg := <-result
	assert.NoError(t, msg.Error)
}

func Test_voidStorageSaveChanged(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileB := &core.File{
		Path:     testFilePath,
		Checksum: "abX",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), TargetName, fileB).Return(nil)

	store := &voidStorage{
		fileStore: fileStore,
	}

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, fileB, nil)
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: true,
		Error:   nil,
		Message: "requested operation was successful",
	}, data)
}

func Test_voidStore_Delete(t *testing.T) {
	store := &voidStorage{}
	result := make(chan core.TargetOperationResult)
	go store.Delete(context.TODO(), result, []*core.File{}, &core.TargetStorageDeleteOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   core.ErrNotImplemented,
		Message: "not implemented",
	}, data)
}

func Test_voidStore_Restore(t *testing.T) {
	store := &voidStorage{}
	result := make(chan core.TargetOperationResult)
	go store.Restore(context.TODO(), result, []*core.File{}, &core.TargetStorageRestoreOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   core.ErrNotImplemented,
		Message: "not implemented",
	}, data)
}

func Test_voidStore_Info(t *testing.T) {
	store := &voidStorage{}
	info := store.Info(context.TODO())
	assert.NotNil(t, info)
}

func Test_voidStore_Ping(t *testing.T) {
	store := &voidStorage{}
	err := store.Ping(context.TODO())
	assert.NoError(t, err)
}
