package keyvalue_db

import (
	"github.com/go-bluestore/bluestore/kv/rocksdb_store"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
)

//type KeyValueDB struct {
//}

type KeyValueDB interface {
	SetMergeOperator(string)

	SetCacheSize(uint64)

	Init(string) error

	CreateAndOpen(string) error

	Open(string) error
}

func CreateKeyValueDB(cct *types.CephContext, _type string, dir string, p interface{}) KeyValueDB {
	if _type == "rocksdb" {
		return rocksdb_store.CreateRocksDBStore(cct, dir, p)
	} else {
		log.Error("only support rocksdb now.")
	}
	return nil
}
