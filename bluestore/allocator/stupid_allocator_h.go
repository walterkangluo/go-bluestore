package allocator

import (
	"github.com/go-bluestore/bluestore/types"
	"sync"
)

type allocatorT map[uint64]uint64
type intervalSetMapT map[uint64]allocatorT
type intervalSetT map[uint64]intervalSetMapT
type StupidAllocator struct {
	name      string
	cct       *types.CephContext
	lock      sync.Mutex
	numFree   int64
	free      []intervalSetT
	lastAlloc uint64
}
