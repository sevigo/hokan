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

func (c *creator) Create(ctx context.Context, data *core.EventData) error {
	return nil
}
