package testnotify

import (
	notifycore "github.com/sevigo/notify/core"
	"github.com/sevigo/notify/event"

	"github.com/sevigo/hokan/pkg/core"
)

type FakeWatcher struct {
	events chan event.Event
	errors chan event.Error
}

func New() core.Notifier {
	return &FakeWatcher{
		events: make(chan event.Event),
		errors: make(chan event.Error),
	}
}

func (w *FakeWatcher) Event() chan event.Event {
	return w.events
}
func (w *FakeWatcher) Error() chan event.Error {
	return w.errors
}

func (w *FakeWatcher) RescanAll() {}

func (w *FakeWatcher) StopWatching(path string) {}

func (w *FakeWatcher) StartWatching(path string, _ *notifycore.WatchingOptions) {
	w.events <- event.Event{
		Action: event.FileAdded,
		Path:   path,
	}
}
