package watcher

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/event"
)

func TestWatch_StartDirWatcher(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	wg.Add(1)
	e := event.New(event.Config{})

	notifier := mocks.NewMockNotifier(controller)
	notifier.EXPECT().StartWatching("/foo", gomock.Any())

	ctx := context.Background()
	w := &Watch{
		ctx:      ctx,
		event:    e,
		notifier: notifier,
	}

	go w.StartDirWatcher()
	wg.Wait()

	w.event.Publish(ctx, &core.EventData{
		Type: core.WatchDirStart,
		Data: &core.Directory{
			Active: true,
			Path:   "/foo",
		},
	})

	time.Sleep(100 * time.Millisecond)
	ctx.Done()
}

func TestWatch_StartRescanWatcher(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	wg.Add(1)
	e := event.New(event.Config{})

	notifier := mocks.NewMockNotifier(controller)
	notifier.EXPECT().RescanAll()

	ctx := context.Background()
	w := &Watch{
		ctx:      ctx,
		event:    e,
		notifier: notifier,
	}

	go w.StartRescanWatcher()
	wg.Wait()

	w.event.Publish(ctx, &core.EventData{
		Type: core.WatchDirRescan,
	})

	time.Sleep(100 * time.Millisecond)
	ctx.Done()
}

func TestWatch_GetDirsToWatch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mocks.NewMockDirectoryStore(controller)
	store.EXPECT().List(gomock.Any()).Return([]*core.Directory{
		{
			ID:     "1",
			Active: true,
			Path:   "/foo",
		},
		{
			ID:     "2",
			Active: false,
			Path:   "/bar",
		},
	}, nil)
	notifier := mocks.NewMockNotifier(controller)
	e := mocks.NewMockEventCreator(controller)
	e.EXPECT().Publish(gomock.Any(), &core.EventData{
		Type: core.WatchDirStart,
		Data: &core.Directory{
			ID:     "1",
			Active: true,
			Path:   "/foo",
		},
	}).Return(nil)

	ctx := context.Background()
	w := &Watch{
		event:    e,
		ctx:      ctx,
		store:    store,
		notifier: notifier,
	}
	w.GetDirsToWatch()
}
