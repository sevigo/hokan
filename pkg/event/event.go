package event

import (
	"context"

	"github.com/sevigo/hokan/pkg/core"
)

type Config struct{}

type creator struct{}

func New(config Config) core.EventCreator {
	return creator{}
}

// Create(context.Context, *EventData) error
func (c creator) Publish(ctx context.Context, data *core.EventData) error {
	return nil
}
