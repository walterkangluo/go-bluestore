package rocksdb_store

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"github.com/go-bluestore/bluestore/types"
	lrdb "github.com/go-bluestore/lib/gorocksdb"
	"github.com/go-bluestore/log"
	"sync"
	"syscall"
)

const (
	lRocksDBFirst = 34300 + iota
	lRocksDBGets
	lRocksDBTxns
	lRocksDBTxnsSync
	lRocksDBGetLatency
	lRocksDBSubmitLatency
	lRocksDBSubmitSyncLatency
	lRocksDBCompact
	lRocksDBCompactRange
	lRocksDBCompactQueueMerge
	lRocksDBCompactQueueLen
	lRocksDBWriteWalTime
	lRocksDBWriteMemTableTime
	lRocksDBWriteDelayTime
	lRocksDBWritePreAndPostProcessTime
	lRocksDBLast
)

type RocksDBStore struct {
	cct    *types.CephContext
	logger *types.PerfCounters
	path   string
	priv   interface{}

	DB      *lrdb.DB
	Env     *BlueRocksEnv
	BbtOpts *lrdb.BlockBasedTableOptions

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
	rs.Env = p.(*BlueRocksEnv)
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

func (rs *RocksDBStore) CreateAndOpen(stream string) error {
	var r error
	if nil != rs.Env {
		var result BlueRocksDirectory
		r = rs.Env.NewDirectory(stream, &result)
		if r != nil {
			log.Error("failed to create dir %s.", stream)
			return r
		}
	} else {
		r = syscall.Mkdir(rs.path, 0755)
		if r != nil || r != syscall.EEXIST {
			log.Error("failed to create %s.", rs.path)
			return r
		}
	}

	return rs.doOpen(stream, true)
}

func (rs *RocksDBStore) Open(stream string) error {
	return rs.doOpen(stream, false)
}
