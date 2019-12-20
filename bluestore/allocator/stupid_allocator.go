package allocator

import (
	btypes "github.com/go-bluestore/bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/types"
)

func CreateStupidAllocator(cct *types.CephContext, name string) *StupidAllocator {
	return &StupidAllocator{
		name: name,
		cct:  cct,
	}
}

func (sa *StupidAllocator) Allocate(wantSize uint64, allocUnit uint64,
	maxAllocSize uint64, hint int64, extents *btypes.PExtentVector) int64 {

	return int64(0)
}

func (sa *StupidAllocator) AllocateInit(wantSize uint64, allocUnit uint64,
	hint uint64, offset *uint64, length *uint64) int64 {
	return int64(0)
}

func (sa *StupidAllocator) Release(releaseSet []uint64) {
	return
}

func (sa *StupidAllocator) GetFree() uint64 {
	return uint64(0)
}

func (sa *StupidAllocator) InitAddFree(offset uint64, length uint64) {
	return
}

func (sa *StupidAllocator) InitRmFree(offset uint64, length uint64) {
	return
}

func (sa *StupidAllocator) Shutdown(offset uint64, length uint64) {
	return
}

func (sa *StupidAllocator) Dump() {
	return
}
