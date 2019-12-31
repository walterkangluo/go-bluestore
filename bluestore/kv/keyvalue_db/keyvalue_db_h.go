package keyvalue_db

import (
	"github.com/go-bluestore/bluestore/kv/rocksdb_store"
	"github.com/go-bluestore/bluestore/types"
	"unsafe"
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

func CreateKeyValueDB(cct *types.CephContext, _type string, dir string, p unsafe.Pointer) KeyValueDB {
	if _type == "leveldb" {
		return nil
	}

	if _type == "rocksdb" {
		return rocksdb_store.CreateRocksDBStore(cct, dir, p)
	}

	if _type == "memdb" {
		return nil
	}

	return nil
}
