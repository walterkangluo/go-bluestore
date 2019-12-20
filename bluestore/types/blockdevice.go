package types

import (
	"github.com/go-bluestore/lib/aio"
	"sync"
	"unsafe"
)

type IOContext struct {
	lock          sync.Mutex
	conditionCond sync.Cond
	r             int

	Cct           *CephContext
	Priv          unsafe.Pointer
	NvmeTaskFirst unsafe.Pointer
	NvmeTaskLast  unsafe.Pointer

	PendingAios []aio.AIOT
	RunningAios []aio.AIOT
	NumPending  int
	NumRunning  int
	AllowEio    bool
}

func CreateIOContext(cct *CephContext, p unsafe.Pointer, _allowAio bool) *IOContext {
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

type BlockDevice struct {
	cct  *CephContext
	path string
}

func CreateBlockDevice(cct *CephContext, path string) *BlockDevice {
	return &BlockDevice{
		cct:  cct,
		path: path,
	}
}

func (*BlockDevice) Open(path string) {

}

func (bd *BlockDevice) GetSize() uint64 {
	return uint64(1)
}
