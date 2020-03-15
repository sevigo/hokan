package main

import (
	"context"

	"github.com/google/wire"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target"
)

var targetSet = wire.NewSet(
	provideTarget,
)

func provideTarget(ctx context.Context, fileStore core.FileStore, event core.EventCreator) (core.TargetRegister, error) {
	return target.New(ctx, fileStore, event)
}
