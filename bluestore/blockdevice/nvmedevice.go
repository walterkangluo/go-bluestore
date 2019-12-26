package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	"unsafe"
)

const (
	dataBufferDefaultNum = 1024
	dataBufferSize       = 8192
	inlineSegmentNum     = 21
	queueId              = -1
)

type sharedDriverQueueDriver struct {
	driver *sharedDriverData
}

type sharedDriverData struct {
	id          uint
	coreId      uint32
	sn          string
	blockSize   uint64
	sectorSize  uint32
	size        uint64
	queueNumber uint32
	queues      []*sharedDriverQueueDriver
}

func CreateNVMEDevice(cct *types.CephContext, path string, cb AioCallbackT, cbPriv unsafe.Pointer) *NVMEDevice {
	nd := &NVMEDevice{
		driver:          nil,
		size:            0,
		blockSize:       0,
		aioStop:         false,
		aioCallback:     cb,
		aioCallbackPriv: cbPriv,
	}
	nd.bufferLock.New("NVMEDevice::buffer_lock")
	nd.BlockDevice.New(cct)

	return nd
}

type IOSegment struct {
	len  uint32
	addr unsafe.Pointer
}

type IORequest struct {
	curSegIdx  uint16
	nseg       uint16
	curSegLeft uint32
	inlineSegs [inlineSegmentNum]unsafe.Pointer
	extraSegs  *unsafe.Pointer
}

type SharedDriverQueueData struct {
	driver *sharedDriverData
}
