package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/lib/aio"
	"sync"
	"unsafe"
)

type AioCallbackT func(handler unsafe.Pointer, aio unsafe.Pointer)

type IOContext struct {
	lock          sync.Mutex
	conditionCond sync.Cond
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

type BlockDevice struct {
	// public
	Cct *types.CephContext

	Path string

	// private
	iocReapLock  sync.Mutex
	iocReapQueue []*IOContext
	iocReapCount int
	rotational   bool
}

func (bd *BlockDevice) IsRotational() bool {
	return bd.rotational
}

func (bd *BlockDevice) New(cct *types.CephContext) {
	bd.Cct = cct
}

func (*BlockDevice) Open(path string) error {
	return nil
}

func (bd *BlockDevice) GetSize() uint64 {
	return uint64(1)
}

func (bd *BlockDevice) GetBlockSize() uint64 {
	return uint64(1)
}

func (bd *BlockDevice) Write(off uint64, bl types.BufferList, buffered bool) {
}

func (bd *BlockDevice) Flush() {
}

func (bd *BlockDevice) Close() {
}

func (bd *BlockDevice) QueueReapIoc() {
}

func (bd *BlockDevice) SupportedBdevLable() bool {
	return true
}
