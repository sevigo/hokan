package event

import (
	"context"
	"sync"

	"github.com/sevigo/hokan/pkg/core"
)

const maxChanSize = 100

type Config struct{}

type creator struct {
	sync.Mutex

	subs map[core.EventType][]chan *core.EventData
}

func New(config Config) core.EventCreator {
	return &creator{
		subs: make(map[core.EventType][]chan *core.EventData),
	}
}

// Create(context.Context, *EventData) error
func (c *creator) Publish(ctx context.Context, data *core.EventData) error {
	c.Lock()
	defer c.Unlock()

	eventType := data.Type
	for _, c := range c.subs[eventType] {
		c <- data
	}

	return nil
}

// Subscribe to a specific event
func (c *creator) Subscribe(ctx context.Context, eventType core.EventType) <-chan *core.EventData {
	c.Lock()
	defer c.Unlock()
	handler := make(chan *core.EventData, maxChanSize)
	_, ok := c.subs[eventType]
	if !ok {
		c.subs[eventType] = []chan *core.EventData{}
	}
	c.subs[eventType] = append(c.subs[eventType], handler)

	go func() {
		<-ctx.Done()
		c.Lock()
		defer c.Unlock()
		delete(c.subs[eventType], handler)
		close(handler)
	}()

	return handler
}
