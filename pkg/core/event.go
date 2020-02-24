package core

import "context"

type EventType int

const (
	WatchDirStart EventType = iota
	WatchDirStop
	FileAdded
	FileChanged
	FileRemoved
)

type EventData struct {
	Type EventType
	Data []byte
}

type EventCreator interface {
	Publish(context.Context, *EventData) error
	Subscribe(context.Context, EventType) <-chan *EventData
}
