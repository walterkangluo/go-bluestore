package rocksdb_store

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/lib/rockdb"
	"sync"
)

type RocksDBStore struct {
	cct    *types.CephContext
	logger *types.PerfCounters
	path   string
	priv   interface{}

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

func CreateRocksDBStore(c *types.CephContext, path string, p interface{}) (rs *RocksDBStore) {
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
	rs.EnableRmRange = c.Conf.RocksDBEnableRmRange
	rs.HighPriWatermark = 0
	return
}

func (rs *RocksDBStore) SetMergeOperator(string) {
}

func (rs *RocksDBStore) SetCacheSize(uint64) {
}

func (rs *RocksDBStore) Init(string) error {
	return nil
}

func (rs *RocksDBStore) CreateAndOpen(str string) error {
	return nil
}

func (rs *RocksDBStore) Open(str string) error {
	return rs.doOpen(str, false)
}
