package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/lib/aio"
	"sync"
	"unsafe"
)

const RW_IO_MAX  = 0x7FFFF000

type AIOCompletionThread struct {
	bdev *BlockDevice
}

func (aio *AIOCompletionThread) New(bdev *BlockDevice) {
	aio.bdev = bdev
}

type KernelDevice struct {
	*BlockDevice
	fdDirect   int
	fdBuffered int
	size       uint64
	blockSize  uint64
	path       string
	fs         types.FS
	aio        bool
	dio        bool

	debugLock     types.Mutex
	debugInflight []uint64

	ioSinceFlush bool
	flushMutex   sync.Mutex

	aioQueue        aio.AioQueueT
	aioCallback     AioCallbackT
	aioCallbackPriv unsafe.Pointer
	aioStop         bool
	injectingCrash  int
	aioThread       AIOCompletionThread
}

func (kr *KernelDevice) Read() {
	kr.BlockDevice.re
}