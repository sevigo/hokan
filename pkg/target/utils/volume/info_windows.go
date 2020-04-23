// +build windows

package volume

// #include "info_windows.h"
import "C"
import (
	"context"
	"fmt"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

type res struct {
	total uint64
	free  uint64
}

var resChan chan res

func init() {
	resChan = make(chan res, 1)
}

func GetVolumeInformation(ctx context.Context, value string) (uint64, uint64) {
	cvalue := C.CString(value)
	fmt.Printf(">>> value=%s\n", value)

	defer func() {
		C.free(unsafe.Pointer(cvalue))
	}()
	go C.GetVolumeInfo(cvalue)

	select {
	case <-ctx.Done():
		log.Info("GetVolumeInformation(): event stream canceled")
		return 0, 0
	case d := <-resChan:
		return d.free, d.total
	case <-time.After(1 * time.Second):
		close(resChan)
		return 0, 0
	}
}

//export goCallbackVolumeInformation
func goCallbackVolumeInformation(freeC, totalC C.longlong) {
	total := uint64(totalC) / 1024 / 1024
	free := uint64(freeC) / 1024 / 1024
	resChan <- res{
		total: total,
		free:  free,
	}
}
