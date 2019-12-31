package rocksdb_store

import (
	"github.com/go-bluestore/bluestore"
	"github.com/go-bluestore/bluestore/types"
	lrdb "github.com/go-bluestore/lib/rockdb"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"strings"
	"syscall"
)

type dbPath struct {
	path string
	size string
}

func parsePath(path string) []dbPath {
	split := strings.Split(path, ",")
	if len(split)%2 != 0 {
		log.Error("invalid path %s.", path)
	}

	dbPaths := make([]dbPath, 0)
	for i := 0; i < len(split); i = i + 2 {
		dbPaths = append(dbPaths, dbPath{path: split[i], size: split[i+1]})
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

type RockDBOpt struct {
	// confirm return type of CreateDBStatistics()
	statistics      string
	walDir          string
	createIfMissing bool
	dbPaths         []dbPath
	infoLog         *cephRocksDBLogger
	env             *lrdb.Env
}

// TODO: parse options
func parseOptionsFromString(optStr string, opt *RockDBOpt) error {
	//strMap := make(map[string]string)

	return nil
}

func (rs *RocksDBStore) doOpen(str string, createIfMissing bool) error {
	var opt RockDBOpt
	if 0 != len(rs.optionStr) {
		r := parseOptionsFromString(rs.optionStr, &opt)
		if nil != r {
			return syscall.EINVAL
		}
	}

	if bluestore.GConf.RocksDBPerf {
		//opt.statistics = dbstats
	}

	opt.createIfMissing = createIfMissing
	if bluestore.GConf.RocksDBSeperateWalDir {
		opt.walDir = rs.path + ".wal"
	}

	// 解析DB path信息，格式为： path：size
	if 0 != len(bluestore.GConf.RocksDBPaths) {
		opt.dbPaths = parsePath(bluestore.GConf.RocksDBPaths)
		utils.AssertTrue(0 != len(opt.dbPaths))
	}

	if bluestore.GConf.RocksDBLogToCephLog {
		opt.infoLog = nil
		opt.infoLog = createCephRocksDBLogger(bluestore.GCephContext)
	}

	if rs.priv != nil {
		log.Debug("using ceph env %v", rs.priv.(*lrdb.Env))
		opt.env = rs.priv.(*lrdb.Env)
	}

	// caches
	if rs.setCacheFlag {
		rs.cacheSize = bluestore.GConf.BlueStoreCacheSize
	}

	//rowCacheSize := uint64(float64(rs.cacheSize) * bluestore.GConf.RocksDBCacheRowRatio)
	//blockSize := rs.cacheSize - rowCacheSize

	return nil
}
