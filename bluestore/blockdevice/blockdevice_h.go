package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	types2 "github.com/go-bluestore/common/types"
	"github.com/go-bluestore/lib/aio"
	"github.com/go-bluestore/utils"
	"sync"
	"sync/atomic"
	"unsafe"
)

type AioCallbackT func(handler unsafe.Pointer, aio unsafe.Pointer)

type IOContext struct {
	lock          sync.Mutex
	conditionCond *sync.Cond
	r             int

	Cct           *types.CephContext
	Priv          unsafe.Pointer
	NvmeTaskFirst unsafe.Pointer
	NvmeTaskLast  unsafe.Pointer

	PendingAios []aio.AIOT
	RunningAios []aio.AIOT
	NumPending  int
	NumRunning  int
	AllowEio    bool
}

func CreateIOContext(cct *types.CephContext, p unsafe.Pointer, _allowAio bool) *IOContext {
	return &IOContext{
		Cct:      cct,
		Priv:     p,
		AllowEio: _allowAio,
	}
}

func (io *IOContext) HasPendingAios() int {
	return io.NumPending
}

func (io *IOContext) SetReturnValue(_r int) {
	io.r = _r
}

func (io *IOContext) GetReturnValue(_r int) int {
	return io.r
}

func (io *IOContext) AioWait() {
}

func (io *IOContext) TryAioAwake() {
	io.conditionCond = sync.NewCond(&io.lock)
	if io.NumRunning == 1 {
		io.conditionCond.Broadcast()
		io.NumRunning--
		utils.AssertTrue(io.NumRunning >= 0)
	} else {
		io.NumRunning--
	}
}

type BlockDevice struct {
	// public
	Cct *types.CephContext

	Path string

	// private
	iocReapLock  sync.Mutex
	iocReapQueue []*IOContext
	iocReapCount atomic.Value //should use atomit Int32
	rotational   bool

	// virtual function
	BlockDeviceFunc interface {
		SupportedBdevLable() bool
		IsRotational() bool
		AioSubmit(ioc *IOContext)
		GetSize() uint64
		GetBlockSize() uint64
		CollectMetadata(prefix string, pm *map[string]string) error
		Read(off uint64, len uint64, pbl *types.BufferList, ioc *IOContext, buffered bool) error
		ReadRandom(off uint64, len uint64, buf string, buffered bool) error
		Write(off uint64, bl *types.BufferList, buffered bool) error
		AioRead(off uint64, len uint64, pbl *types.BufferList, ioc *IOContext) error
		AioWrite(off uint64, bl *types2.List, ioc *IOContext, buffered bool) bool
		Flush() error
		InvalidateCache(off uint64, len uint64) error
		Open(path string) error
		Close()
    }
}

func (bd *BlockDevice) New(cct *types.CephContext) {
	bd.Cct = cct
}

func (bd *BlockDevice) QueueReapIoc() {
}
