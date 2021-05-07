package void

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

var testFilePath = "/test/file.txt"
var testErr = errors.New("test error")

func Test_voidStorageSaveError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	file := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), name, file).Return(testErr)

	store := &voidBackup{
		fileStore: fileStore,
	}

	result := make(chan core.BackupResult)
	go store.Save(context.TODO(), result, file, nil)
	msg := <-result
	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   testErr,
	}, msg)
}

func Test_voidStorageSaveOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	file := &core.File{
		Path:     testFilePath,
		Checksum: "abX",
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), name, file).Return(nil)

	store := &voidBackup{
		fileStore: fileStore,
	}

	result := make(chan core.BackupResult)
	go store.Save(context.TODO(), result, file, nil)
	msg := <-result
	assert.Equal(t, core.BackupResult{
		Success: true,
		Error:   nil,
		Message: core.BackupSuccessMessage,
	}, msg)
}

func Test_voidStore_Delete(t *testing.T) {
	store := &voidBackup{}
	result := make(chan core.BackupResult)
	go store.Delete(context.TODO(), result, []*core.File{}, &core.BackupOperationOptions{})
	data := <-result
	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}, data)
}

func Test_voidStore_Restore(t *testing.T) {
	store := &voidBackup{}
	result := make(chan core.BackupResult)
	go store.Restore(context.TODO(), result, []*core.File{}, &core.BackupOperationOptions{})
	data := <-result
	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}, data)
}

func Test_voidStore_Ping(t *testing.T) {
	store := &voidBackup{}
	err := store.Ping(context.TODO())
	assert.NoError(t, err)
}

func Test_voidBackup_Name(t *testing.T) {
	store := &voidBackup{}
	name := store.Name()
	assert.Equal(t, "void", name)
}
