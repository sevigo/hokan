package target

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/minio"
	"github.com/sevigo/hokan/pkg/target/void"
)

// All known targets
var targets = map[string]core.TargetFactory{
	void.TargetName:  void.New,
	minio.TargetName: minio.New,
}

type Register struct {
	ctx context.Context
	sync.Mutex
	fileStore core.FileStore
	event     core.EventCreator
	register  map[string]core.TargetStorage
}

func New(ctx context.Context, fileStore core.FileStore, event core.EventCreator) (*Register, error) {
	log.Debug().Str("op", "target.New()").Msg("start")
	r := &Register{
		ctx:       ctx,
		fileStore: fileStore,
		event:     event,
		register:  make(map[string]core.TargetStorage),
	}
	r.InitTargets(ctx)
	go r.StartFileAddedWatcher()
	return r, nil
}

func (r *Register) InitTargets(ctx context.Context) {
	for name, target := range targets {
		go r.initWithRetry(ctx, name, target)
	}
}

func (r *Register) initWithRetry(ctx context.Context, name string, target core.TargetFactory) {
	counter := 0
	mod := 10
	ticker := time.NewTicker(time.Duration(mod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Str("op", "target.retryTarger()").Str("target", name).Msg("stream canceled")
			return
		case <-ticker.C:
			log.Debug().
				Str("op", "target.retryTarger()").
				Str("target", name).Str("offset", fmt.Sprintf("%d sec", mod)).
				Str("counter", fmt.Sprintf("%d", counter)).
				Msg("retry target setup")
			ts, err := target(ctx, r.fileStore)
			if err == nil {
				r.addTarget(name, ts)
				if counter > 0 {
					fmt.Println(">>> send event WatchDirRescan")
					r.event.Publish(ctx, &core.EventData{
						Type: core.WatchDirRescan,
					})
				}
				return
			}
			log.Error().Err(err).Str("target", name).Msg("Can't create new target storage")
			if counter%10 == 0 {
				mod = mod * 2
				ticker = time.NewTicker(time.Duration(mod) * time.Second)
			}
			counter++
		}
	}
}

func (r *Register) StartFileAddedWatcher() {
	ctx := r.ctx
	eventData := r.event.Subscribe(ctx, core.FileAdded)

	for {
		select {
		case <-ctx.Done():
			log.Printf("file-watcher: stream canceled")
			return
		case e := <-eventData:
			log.Debug().Msgf("StartFileAddedWatcher(): %#v\n", e.Data)
			err := r.processFileAddedEvent(e)
			if err != nil {
				log.Err(err).Msg("can't send the file to the target storage")
			}
		}
	}
}

func (r *Register) processFileAddedEvent(e *core.EventData) error {
	file, ok := e.Data.(core.File)
	if !ok {
		return fmt.Errorf("invalid event data: %v", e)
	}
	for _, target := range file.Targets {
		if ts := r.getTarget(target); ts != nil {
			err := ts.Save(r.ctx, &file)
			if err != nil {
				log.Err(err).
					Str("target", target).
					Str("file", file.Path).
					Msg("can't save the file to the target storage")
			}
		}
	}
	return nil
}

func (r *Register) addTarget(name string, ts core.TargetStorage) {
	r.Lock()
	defer r.Unlock()
	r.register[name] = ts
}

func (r *Register) getTarget(name string) core.TargetStorage {
	r.Lock()
	defer r.Unlock()
	return r.register[name]
}
