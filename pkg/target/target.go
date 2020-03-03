package target

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/void"
)

// All known targets
var targets = map[string]core.TargetFactory{
	void.TargetName: void.New,
}

var registerM sync.Mutex
var register map[string]core.TargetStorage = make(map[string]core.TargetStorage)

// type Target struct {
// 	ctx      context.Context
// 	event    core.EventCreator
// 	store    core.DirectoryStore
// 	notifier core.Notifier
// 	catalog  []*core.Directory
// }

func New(ctx context.Context, fileStore *core.FileStore, event core.EventCreator) {
	initTargets(fileStore)
}

func initTargets(fileStore *core.FileStore) {
	for name, target := range targets {
		ts, err := target(fileStore)
		if err != nil {
			log.Err(err).Msg("Can't create new target storage")
			continue
		}
		addTarget(name, ts)
	}
}

func addTarget(name string, ts core.TargetStorage) {
	registerM.Lock()
	defer registerM.Unlock()
	register[name] = ts
}
