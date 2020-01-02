package rocksdb_store

import (
	"github.com/go-bluestore/bluestore"
	"github.com/go-bluestore/bluestore/types"
	lrdb "github.com/go-bluestore/lib/gorocksdb"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"strconv"
	"strings"
	"syscall"
)

func parsePath(path string) []*lrdb.DBPath {
	split := strings.Split(path, ",")
	if len(split)%2 != 0 {
		log.Error("invalid path %s.", path)
	}

	m := len(split) / 2
	dbPaths := make([]*lrdb.DBPath, m)
	for i := 0; i < len(split); i = i + 2 {
		path := split[i]
		size, err := strconv.ParseInt(split[i+1], 10, 64)
		utils.AssertTrue(err == nil)
		dbPaths[i] = lrdb.NewDBPath(path, uint64(size))
	}
	return dbPaths
}

type cephRocksDBLogger struct {
	cct *types.CephContext
}

func createCephRocksDBLogger(cct *types.CephContext) *cephRocksDBLogger {
	return &cephRocksDBLogger{
		cct: cct,
	}
}

type MergeOperatorRouter struct {
	store *RocksDBStore
}

func NewMergeOperatorRouter(s *RocksDBStore) *MergeOperatorRouter {
	return &MergeOperatorRouter{store: s}
}

// TODO:
func (mr *MergeOperatorRouter) Name() string {
	return ""
}

// TODO:
func (mr *MergeOperatorRouter) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	return nil, false
}

func (rs *RocksDBStore) compact() {
	rs.logger.Inc(lRocksDBCompact, uint64(1))
	rs.DB.CompactRange(lrdb.Range{Start: nil, Limit: nil})
}

func (rs *RocksDBStore) doOpen(str string, createIfMissing bool) error {
	var r error
	opts := lrdb.NewDefaultOptions()
	opts, r = lrdb.GetOptionsFromString(opts, rs.optionStr)
	if nil != r {
		log.Error("error option string %s.", rs.optionStr)
		return syscall.EINVAL
	}

	if bluestore.GConf.RocksDBPerf {
		opts.EnableStatistics()
	}

	opts.SetCreateIfMissing(createIfMissing)
	if bluestore.GConf.RocksDBSeperateWalDir {
		opts.SetWalDir(rs.path + ".wal")
	}

	// path：size
	if 0 != len(bluestore.GConf.RocksDBPaths) {
		opts.SetDBPaths(parsePath(bluestore.GConf.RocksDBPaths))
	}

	if bluestore.GConf.RocksDBLogToCephLog {
		opts.SetInfoLogLevel(lrdb.InfoInfoLogLevel)
	}

	if rs.priv != nil {
		log.Debug("using ceph env %v", rs.priv.(*lrdb.Env))
		opts.SetEnv(rs.priv.(*lrdb.Env))
	}

	// caches
	if rs.setCacheFlag {
		rs.cacheSize = bluestore.GConf.BlueStoreCacheSize
	}

	rowCacheSize := uint64(float64(rs.cacheSize) * bluestore.GConf.RocksDBCacheRowRatio)
	blockCacheSize := rs.cacheSize - rowCacheSize

	if bluestore.GConf.RocksDBCacheType == "lru" {
		rs.BbtOpts.SetBlockCache(lrdb.NewLRUCache(uint64(blockCacheSize)))
	} else {
		log.Error("not support cache type %s now.", bluestore.GConf.RocksDBCacheType)
		return syscall.EINVAL
	}
	rs.BbtOpts.SetBlockSize(bluestore.GConf.RocksDBBlockSize)

	// TODO: confirm opt row cache
	//if rowCacheSize > 0 {
	//	opt.rowCache = lrdb.NewLRUCache(rowCacheSize)
	//	opts.
	//}

	bloomBits := bluestore.GConf.RocksDBBloomBitsPerKey
	if bloomBits > 0 {
		rs.BbtOpts.SetFilterPolicy(lrdb.NewBloomFilter(int(bloomBits)))
	}

	if "binary_search" == bluestore.GConf.RocksDBIndexType {
		rs.BbtOpts.SetIndexType(lrdb.KBinarySearchIndexType)
	}

	if "hash_search" == bluestore.GConf.RocksDBIndexType {
		rs.BbtOpts.SetIndexType(lrdb.KHashSearchIndexType)
	}

	if "two_search" == bluestore.GConf.RocksDBIndexType {
		rs.BbtOpts.SetIndexType(lrdb.KTwoLevelIndexSearchIndexType)
	}

	rs.BbtOpts.SetCacheIndexAndFilterBlocks(bluestore.GConf.RocksDBCacheIndexAndFilterBlocks)
	rs.BbtOpts.SetCacheIndexAndFilterBlocksWithHighPriority(bluestore.GConf.RocksDBCacheIndexAndFilterBlocksWithHighProority)
	rs.BbtOpts.SetPartitionFilters(bluestore.GConf.RocksDBPartitionFilters)
	rs.BbtOpts.SetMetadataBlockSize(bluestore.GConf.RockdSBMetadataBlockSize)
	rs.BbtOpts.SetPinL0FilterAndIndexBlocksInCache(bluestore.GConf.RocksDBPinL0FilterAndIndexBlocksInCache)

	opts.SetBlockBasedTableFactory(rs.BbtOpts)

	log.Debug("block_size %d, block_cache_size %d, row_cache_size %d.",
		bluestore.GConf.RocksDBBlockSize, blockCacheSize, rowCacheSize)

	opts.SetMergeOperator(NewMergeOperatorRouter(rs))
	var err error
	rs.DB, err = lrdb.OpenDb(opts, rs.path)
	if err != nil {
		log.Error("open db %s failed with err %v.", rs.path, err)
		return syscall.EINVAL
	}

	plb := types.CreatePerCountersBuilder(bluestore.GCephContext, "rocksdb", lRocksDBFirst, lRocksDBLast)

	plb.AddU64Counter(lRocksDBGets, "get", "Gets", "", 0, 0)
	plb.AddU64Counter(lRocksDBTxns, "submit_transaction", "Submit transaction", "", 0, 0)
	plb.AddU64Counter(lRocksDBTxnsSync, "submit_transaction_sync", "Submit transaction sync", "", 0, 0)

	plb.AddTimeAvg(lRocksDBGetLatency, "get_latency", "Get latency", "", 0, 0)
	plb.AddTimeAvg(lRocksDBSubmitLatency, "submit_latency", "Submit latency", "", 0, 0)
	plb.AddTimeAvg(lRocksDBSubmitSyncLatency, "submit_sync_latency", "Submit latency sync", "", 0, 0)

	plb.AddU64Counter(lRocksDBCompact, "compact", "Compactions", "", 0, 0)
	plb.AddU64Counter(lRocksDBCompactRange, "compact_range", "Compactions by range", "", 0, 0)
	plb.AddU64Counter(lRocksDBCompactQueueMerge, "compact_queue_merge", "Mergings of ranges in compaction queue", "", 0, 0)
	plb.AddU64Counter(lRocksDBCompactQueueLen, "compact_queue_len", "Length of compaction queue", "", 0, 0)

	plb.AddTimeAvg(lRocksDBWriteWalTime, "rocksdb_write_wal_time", "Rocksdb write wal time", "", 0, 0)
	plb.AddTimeAvg(lRocksDBWriteMemTableTime, "rocksdb_write_memetable_time", "Rocksdb write memtable time", "", 0, 0)
	plb.AddTimeAvg(lRocksDBWriteDelayTime, "rocksdb_write_delay_time", "Rocksdb write delay time", "", 0, 0)
	plb.AddTimeAvg(lRocksDBWritePreAndPostProcessTime, "rocksdb_write_pre_and_post_time", "Total time spent on writing a record, excluding write process", "", 0, 0)

	logger := plb.CreatePerfCounters()
	rs.cct.GetPerfCountersCollection().Add(logger)

	if rs.CompactOnMount {
		log.Error("Compacting rocksdb store...")
		rs.compact()
		log.Error("Finished compacting rocksdb store ")
	}

	return nil
}
