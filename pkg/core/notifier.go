package core

import "github.com/sevigo/notify/event"

// Notifier is an interface for start/stop watching directories for changes
type Notifier interface {
	Event() chan event.Event
	Error() chan event.Error
	Scan(string) error
	StopWatching(string)
	StartWatching(string)
}
