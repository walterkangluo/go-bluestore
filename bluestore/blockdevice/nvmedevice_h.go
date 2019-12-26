package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	"sync"
	"unsafe"
)

var (
	ReadCommand  = 0
	WriteCommand = 1
	FlushCommand = 2
)

type Extent struct {
	xLen    uint64
	xOff    uint32
	data    string
	dataLen uint64
}

type Task struct {
	device     *NVMEDevice
	ctx        *IOContext
	command    int
	offset     uint64
	len        uint64
	writeBl    types.BufferList
	next       *Task
	returnCode int
	ioRequest  IORequest
	lock       sync.Mutex
	cond       sync.Cond
	queue      SharedDriverQueueData
}

type BufferedExtents struct {
	offset          uint64
	bufferedExtents map[uint64]Extent // offset, exyent
	leftEdge        uint64
	rightEdge       uint64
}

type NVMEDevice struct {
	*BlockDevice
	driver    *sharedDriverData
	name      string
	size      uint64
	blockSize uint64
	aioStop   bool

	bufferLock       types.Mutex
	bufferedExtents  BufferedExtents
	bufferedTaskHead *Task

	// public
	aioCallback     AioCallbackT
	aioCallbackPriv unsafe.Pointer
}
