package allocator

import "github.com/go-bluestore/bluestore/bluestore/types"

type Allocator interface {
	Allocate(wantSize uint64, allocUnit uint64, maxAllocSize uint64, hint int64, extents *types.PExtentVector) int64

	AllocateInit(wantSize uint64, allocUnit uint64, hint uint64, offset *uint64, length *uint64) int64

	Release(releaseSet []uint64)

	GetFree() uint64

	InitAddFree(offset uint64, length uint64)

	InitRmFree(offset uint64, length uint64)

	Shutdown(offset uint64, length uint64)

	Dump()
}
