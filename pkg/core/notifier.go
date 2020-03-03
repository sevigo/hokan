package core

import (
	"github.com/sevigo/notify/core"
	"github.com/sevigo/notify/event"
)

type WatchOptions struct {
	Rescan bool
}

// Notifier is an interface for start/stop watching directories for changes
type Notifier interface {
	Event() chan event.Event
	Error() chan event.Error
	Scan(string) error
	StopWatching(string)
	StartWatching(string, *core.WatchingOptions)
}
