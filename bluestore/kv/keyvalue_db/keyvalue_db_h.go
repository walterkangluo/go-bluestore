package keyvalue_db

import (
	"github.com/go-bluestore/bluestore/kv/rocksdb_store"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
)

type TransactionImpl struct {
}

func (ti *TransactionImpl) set1(prefix string, toSet map[string]types.BufferList) {

}

func (ti *TransactionImpl) set2(prefix string, toSetBl *types.BufferList) {
	p := toSetBl.Begin()
	var num uint32
	toSetBl.Decode(utils.Int32ToBytes(num), p)
	for {
		var key string
		//var value types.BufferList
		toSetBl.Decode([]byte(key), p)
		//toSetBl.Decode(value, p)
	}
}

func (ti *TransactionImpl) set3() {

}

type KeyValueDB struct {
	transaction TransactionImpl
}

type KeyValueDBI interface {
	SetMergeOperator(string)

	SetCacheSize(uint64)

	Init(string) error

	CreateAndOpen(string) error

	Open(string) error
}

func CreateKeyValueDB(cct *types.CephContext, _type string, dir string, p interface{}) KeyValueDBI {
	if _type == "rocksdb" {
		return rocksdb_store.CreateRocksDBStore(cct, dir, p)
	} else {
		log.Error("only support rocksdb now.")
	}
	return nil
}
