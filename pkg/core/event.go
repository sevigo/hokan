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
)

type EventData struct {
	Type EventType
	Data interface{}
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
	default:
		return "unknown"
	}
}
