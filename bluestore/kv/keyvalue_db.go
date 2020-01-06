package kv

import (
	"github.com/go-bluestore/bluestore/kv/common"
	"github.com/go-bluestore/bluestore/kv/rocksdb_store"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
)

type KeyValueDB struct {
	keyValueDBInterface

	cacheBytes []int64
	cacheRatio float64
}

type keyValueDBInterface interface {
	SetMergeOperator(prefix string, opera rocksdb_store.MergeOperator)

	SetCacheSize(uint64)

	Init(string) error

	CreateAndOpen(string) error

	Open(string) error

	GetTransaction() common.Transaction
}

func CreateKeyValueDB(cct *types.CephContext, _type string, dir string, p interface{}) *KeyValueDB {
	if _type == "rocksdb" {
		return &KeyValueDB{
			keyValueDBInterface: rocksdb_store.CreateRocksDBStore(cct, dir, p),
			cacheBytes:          make([]int64, types.Last+1),
			cacheRatio:          float64(0),
		}
	} else {
		log.Error("only support rocksdb now.")
	}
	return nil
}
