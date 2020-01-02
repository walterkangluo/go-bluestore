package freeListManager

import (
	"github.com/go-bluestore/bluestore/kv/keyvalue_db"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/utils"
)

type FreelistManage interface {
	// TODO: wait to implement
	Create(size uint64, graunlarity uint64, txn keyvalue_db.KeyValueDBI)
}

type FreelistManager struct {
	Cct *types.CephContext
}

func (fm *FreelistManager) New(cct *types.CephContext) {
	fm.Cct = cct
}

func CreateFreelistManage(cct *types.CephContext, _type string, kvdb keyvalue_db.KeyValueDBI, prefix string) *FreelistManage {
	utils.AssertTrue("B" == prefix)
	if _type == "bitmap" {
		return nil
	}
	return nil
}
