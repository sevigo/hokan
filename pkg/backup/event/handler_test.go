package event

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/backup/event/removed"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/event"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pathToFile := "/path/to/file.txt"
	ctx := context.TODO()
	events := event.New(event.Config{})

	backup := mocks.NewMockBackup(controller)
	backup.EXPECT().Name().Return("test")

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, bucket string, file *core.File) {
			assert.Equal(t, "test", bucket)
			assert.Equal(t, pathToFile, file.Path)
		}).Return(nil)

	results := make(chan core.BackupResult)
	event := &core.EventData{
		Type: core.FileRemoved,
		Data: core.File{
			Path: pathToFile,
		},
	}
	delEvent := events.Subscribe(ctx, core.FileRemoved)
	processor := removed.New(&core.EventHandler{
		Backup:    backup,
		FileStore: fileStore,
		Results:   results,
	})
	assert.Equal(t, "file removed", processor.Name())

	go listener(ctx, delEvent, processor)
	err := events.Publish(ctx, event)
	assert.NoError(t, err)
	msg := <-results
	assert.Equal(t, core.BackupResult{
		Success: true,
		Message: core.BackupFileDeletedMessage,
	}, msg)
}

func TestDeleteErr(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pathToFile := "/path/to/file.txt"
	delErr := errors.New("del error")
	ctx := context.TODO()
	events := event.New(event.Config{})

	backup := mocks.NewMockBackup(controller)
	backup.EXPECT().Name().Return("test")

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, bucket string, file *core.File) {
			assert.Equal(t, "test", bucket)
			assert.Equal(t, pathToFile, file.Path)
		}).Return(delErr)

	results := make(chan core.BackupResult)
	event := &core.EventData{
		Type: core.FileRemoved,
		Data: core.File{
			Path: pathToFile,
		},
	}
	delEvent := events.Subscribe(ctx, core.FileRemoved)
	processor := removed.New(&core.EventHandler{
		Backup:    backup,
		FileStore: fileStore,
		Results:   results,
	})
	assert.Equal(t, "file removed", processor.Name())

	go listener(ctx, delEvent, processor)
	err := events.Publish(ctx, event)
	assert.NoError(t, err)
	msg := <-results
	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   delErr,
	}, msg)
}
