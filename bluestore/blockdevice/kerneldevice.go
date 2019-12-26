package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/lib/aio"
	"unsafe"
)

func CreateKernelDevice(cct *types.CephContext, path string, cb AioCallbackT, cbPriv unsafe.Pointer) *KernelDevice {
	kd := &KernelDevice{
		fdDirect:        -1,
		fdBuffered:      -1,
		size:            0,
		blockSize:       0,
		fs:              types.FS{},
		aio:             false,
		dio:             false,
		aioQueue:        aio.Create(cct.Conf.BdevAioMaxQueueDepth),
		aioCallback:     cb,
		aioCallbackPriv: cbPriv,
		aioStop:         false,
	}
	kd.debugLock.New("kernelDevice::debug_lock")
	kd.BlockDevice.New(cct)
	kd.aioThread.New(kd.BlockDevice)
	return kd
}
