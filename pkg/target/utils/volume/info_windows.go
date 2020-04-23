// +build windows

package volume

// #include "info_windows.h"
import "C"
import (
	"context"
	"unsafe"
)

type res struct {
	total uint64
	free  uint64
}

var resChan chan res

func init() {
	resChan = make(chan res, 1)
}

func GetVolumeInformation(ctx context.Context, path string) (uint64, uint64) {
	cpath := C.CString(path)
	defer func() {
		C.free(unsafe.Pointer(cpath))
	}()
	go C.GetVolumeInfo(cpath)
	d := <-resChan
	return d.free, d.total
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
