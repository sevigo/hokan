package core

import (
	"context"
)

type EventType int

const (
	WatchDirStart EventType = iota
	WatchDirStop
	WatchDirRescan
	FileAdded
	FileChanged
	FileRemoved
	FileRenamed
)

type EventHandler struct {
	Ctx       context.Context
	Events    EventCreator
	Backup    Backup
	FileStore FileStore
	// write results of any operation here, this will be propagated to the UI
	Results chan BackupResult
}

type EventData struct {
	Type EventType
	Data interface{}
}

type EventProcessrFactory func(handler *EventHandler) EventProcessor

type EventProcessor interface {
	Name() string
	Process(event *EventData) error
}

type EventCreator interface {
	Publish(context.Context, *EventData) error
	Subscribe(context.Context, EventType) <-chan *EventData
}

func EventToString(e EventType) string {
	switch e {
	case WatchDirRescan:
		return "dir rescan"
	case WatchDirStop:
		return "watch stop"
	case FileAdded:
		return "file added"
	case FileChanged:
		return "file changed"
	case FileRemoved:
		return "file removed"
	case FileRenamed:
		return "file renamed"
	default:
		return "unknown"
	}
}
