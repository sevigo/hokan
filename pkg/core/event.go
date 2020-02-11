package core

import "context"

type EventData struct{}

type EventCreator interface {
	Create(context.Context, *EventData)
}
