// +build linux

package volume

import (
	"context"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func GetVolumeInformation(ctx context.Context, path string) (uint64, uint64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.WithField("path", path).
			WithError(err).
			Error("linux.GetVolumeInformation(): can't get syscall.Statfs")
		return 0, 0
	}
	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	return uint64(free), uint64(total)
}
