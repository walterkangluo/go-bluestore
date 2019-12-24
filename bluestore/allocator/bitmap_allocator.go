package allocator

import (
	btypes "github.com/go-bluestore/bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/types"
)

func CreateBitmapAllocator(cct *types.CephContext, cap int64, allocUnit int64, name string) *BitmapAllocator {
	return &BitmapAllocator{
		name: name,
		cct:  cct,
	}
}

/*
func (sa *BitmapAllocator) Allocate(wantSize uint64, allocUnit uint64,
	maxAllocSize uint64, hint int64, extents *btypes.PExtentVector) int64 {

	return int64(0)
}
*/

func (sa *BitmapAllocator) Allocate(para ...interface{}) int64 {

	return int64(0)
}

func (sa *BitmapAllocator) AllocateInit(wantSize uint64, allocUnit uint64,
	hint uint64, offset *uint64, length *uint64) int64 {
	return int64(0)
}

func (sa *BitmapAllocator) Release(releaseSet btypes.PExtentVector) {
	return
}

func (sa *BitmapAllocator) GetFree() uint64 {
	return uint64(0)
}

func (sa *BitmapAllocator) InitAddFree(offset uint64, length uint64) {
	return
}

func (sa *BitmapAllocator) InitRmFree(offset uint64, length uint64) {
	return
}

func (sa *BitmapAllocator) Shutdown() {
	return
}

func (sa *BitmapAllocator) Dump() {
	return
}
