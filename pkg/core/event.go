package core

import "context"

type EventData struct{}

type EventCreator interface {
	Publish(context.Context, *EventData) error
}
