package target

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/utils"
)

func (r *Register) StartFileAddedWatcher() {
	ctx := r.ctx
	eventData := r.event.Subscribe(ctx, core.FileAdded)

	for {
		select {
		case <-ctx.Done():
			log.Info("StartFileAddedWatcher(): event stream canceled")
			return
		case e := <-eventData:
			err := r.processFileAddedEvent(e)
			if err != nil {
				log.WithError(err).Error("StartFileAddedWatcher(): can't send the file to the target storage")
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
		if ts := r.GetTarget(target); ts != nil {
			if r.getTargetStatus(target) == core.TargetStoragePaused {
				continue
			}

			// if result.Error != nil {
			// 	log.WithError(result.Error).WithFields(log.Fields{
			// 		"target": target,
			// 		"file":   file.Path,
			// 	}).Error("can't save the file to the target storage")
			// 	r.setTargetStatus(target, core.TargetStorageError)
			// 	continue
			// }
			// if result.Error == nil && r.getTargetStatus(target) == core.TargetStorageError {
			// 	log.WithField("target", target).Info("target storage recoverd")
			// 	r.setTargetStatus(target, core.TargetStorageOK)
			// 	r.rescanAllWatchedDirs()
			// }
		}
	}
	return nil
}

func (r *Register) saveFileToTarget(storage core.TargetStorage, file *core.File) {
	storedFile, err := r.fileStore.Find(r.ctx, &core.FileSearchOptions{
		FilePath:   file.Path,
		TargetName: storage.Name(),
	})
	if errors.Is(err, core.ErrFileNotFound) || utils.FileHasChanged(file, storedFile) {
		storage.Save(r.ctx, r.Results, file, &core.TargetStorageSaveOpt{})
		return
	}
	r.Results <- core.TargetOperationResultSuccess("file hasn't changed")
}

func (r *Register) rescanAllWatchedDirs() {
	log.Info("event.rescanAllWatchedDirs(): publish [WatchDirRescan] event")
	err := r.event.Publish(r.ctx, &core.EventData{Type: core.WatchDirRescan})
	if err != nil {
		log.WithError(err).Error("Can't publish [WatchDirRescan] event")
	}
}
