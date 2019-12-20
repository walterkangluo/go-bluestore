package allocator

import "github.com/go-bluestore/bluestore/types"

type BitmapAllocator struct {
	name string
	cct  *types.CephContext
}
