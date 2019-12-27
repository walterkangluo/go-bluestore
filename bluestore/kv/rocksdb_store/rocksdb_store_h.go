package rocksdb_store

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/lib/rockdb"
	"sync"
	"unsafe"
)

type RocksDBStore struct {
	cct    *types.CephContext
	logger *types.PerfCounters
	path   string
	priv   unsafe.Pointer

	DB  *gorocksdb.DB
	Env *gorocksdb.Env

	//dbStats gorocksdb.Statistics
	//bbtOpts C.BlockBasedTableOptions

	optionStr    string
	cacheSize    uint64
	setCacheFlag bool

	compactQueueLock types.Mutex
	compactQueueCond sync.Cond
	compactQueue     map[string]string
	compactQueueStop bool
	compactThread    CompactThread

	CompactOnMount   bool
	DisableWal       bool
	EnableRmRange    bool
	HighPriWatermark int64
}

type CompactThread struct {
	db *RocksDBStore
}

func (ct *CompactThread) Init(db *RocksDBStore) {
	ct.db = db
}

func (rs *RocksDBStore) New(c *types.CephContext, path string, p unsafe.Pointer) {
	rs.cct = c
	rs.logger = nil
	rs.path = path
	rs.priv = p
	rs.DB = nil
	// TODO: to confirm p
	//rs.Env = gorocksdb.NewNativeEnv(p)
	rs.Env = gorocksdb.NewDefaultEnv()
	//rs.dbStats = nil
	rs.compactQueueLock.New("RocksDBStore::comact_thread_lock")
	rs.compactQueueStop = false
	rs.compactThread.Init(rs)
	rs.CompactOnMount = false
	rs.DisableWal = false
	rs.EnableRmRange = c.Conf.RockDBEnableRmRange
	rs.HighPriWatermark = 0
}
