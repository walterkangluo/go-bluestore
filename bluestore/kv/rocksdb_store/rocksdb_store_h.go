package rocksdb_store

import (
	"github.com/go-bluestore/bluestore/types"
	"sync"
	"unsafe"
)

type RocksDBStore struct {
	cct    *types.CephContext
	logger *types.PerfCounters
	path   string
	priv   unsafe.Pointer

	optionStr    string
	cacheSize    uint64
	setCacheFlag bool

	compactQueueLock types.Mutex
	compactQueueCond sync.Cond
	compactQueue     map[string]string
	compactQueueStop bool
}
