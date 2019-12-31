package allocator

import "github.com/go-bluestore/bluestore/bluestore/types"

type Allocator interface {
	Allocate(para ...interface{}) int64

	AllocateInit(wantSize uint64, allocUnit uint64, hint uint64, offset *uint64, length *uint64) int64

	Release(releaseSet types.PExtentVector)

	GetFree() uint64

	InitAddFree(offset uint64, length uint64)

	InitRmFree(offset uint64, length uint64)

	Shutdown()

	Dump()
}