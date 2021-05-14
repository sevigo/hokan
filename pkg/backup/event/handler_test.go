package event

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/event"
	"github.com/stretchr/testify/assert"
)

func TestInitHanler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.TODO()
	events := event.New(event.Config{})
	backup := mocks.NewMockBackup(controller)
	fileStore := mocks.NewMockFileStore(controller)

	results := make(chan core.BackupResult)

	tests := []struct {
		name   string
		event  *core.EventData
		result core.BackupResult
	}{
		{
			name: "case 1",
			event: &core.EventData{
				Type: core.FileRemoved,
				Data: core.File{
					Path: "/path/to/file.txt",
				},
			},
			result: core.BackupResult{
				Success: true,
				Message: core.BackupFileDeletedMessage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitHanler(ctx, events, backup, fileStore, results)
			events.Publish(ctx, tt.event)
			result := <-results
			assert.Equal(t, tt.result, result)
			// fmt.Printf(">>> %+v\n", result)
		})
	}
	// t.Fail()
}
