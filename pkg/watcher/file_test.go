package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/watcher/testnotify"
)

var testFilePath = "file_test.go"

func TestWatch_StartFileWatcher(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	localPath := filepath.Join(pwd, testFilePath)

	notifier := testnotify.New()
	event := mocks.NewMockEventCreator(controller)
	event.EXPECT().Publish(gomock.Any(), gomock.Any()).Do(func(_ context.Context, e *core.EventData) error {
		data, ok := e.Data.(core.File)
		assert.True(t, ok)
		assert.Equal(t, localPath, data.Path)
		return nil
	})

	ctx := context.Background()
	w := &Watch{
		ctx:      ctx,
		event:    event,
		notifier: notifier,
		catalog: []*core.Directory{
			{
				Path: pwd,
			},
		},
	}

	go w.StartFileWatcher()
	notifier.StartWatching(localPath, nil)
	time.Sleep(100 * time.Millisecond)
	ctx.Done()
}
