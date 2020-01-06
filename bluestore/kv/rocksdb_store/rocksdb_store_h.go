package rocksdb_store

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"github.com/go-bluestore/bluestore/kv/common"
	"github.com/go-bluestore/bluestore/kv/env"
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

type RocksDBTransactionImpl struct {
	common.TransactionImpl
	DB  *RocksDBStore
	Bat lrdb.WriteBatch
}

func NewRocksDBTransactionImpl(db *RocksDBStore) *RocksDBTransactionImpl {
	return &RocksDBTransactionImpl{
		DB: db,
	}
}

type RocksDBStore struct {
	cct    *types.CephContext
	logger *types.PerfCounters
	path   string
	priv   interface{}

	DB      *lrdb.DB
	Env     *env.BlueRocksEnv
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

	mergeOps []mergeOperaPair
}

type mergeOperaPair struct {
	prefix string
	mop    MergeOperator
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
	rs.Env = p.(*env.BlueRocksEnv)
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

type MergeOperator interface {
	MergeNonexistent(rData string, rLen int, newValue *string)

	Merge(rData string, rLen int, lData string, lLen int, newValue *string)

	Name() string
}

func (rs *RocksDBStore) SetCacheSize(uint64) {
}

func (rs *RocksDBStore) Init(string) error {
	return nil
}

func (rs *RocksDBStore) CreateAndOpen(stream string) error {
	var r error
	if nil != rs.Env {
		var result env.BlueRocksDirectory
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

type RocksWBHandler struct {
}
