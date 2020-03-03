package target

import (
	"context"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/void"
)

// All known targets
var targets = map[string]core.TargetFactory{
	void.TargetName: void.New,
}

func New(ctx context.Context, fileStore core.FileStore, event core.EventCreator) {

}

func initTargets(fileStore core.FileStore) {
	for _, target := range targets {

	}
}
