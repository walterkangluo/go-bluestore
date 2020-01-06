package freeListManager

import (
	"github.com/go-bluestore/bluestore/kv"
	"github.com/go-bluestore/bluestore/kv/common"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/utils"
)

type freelistManagerInterface interface {
	// TODO: wait to implement
	Create(size uint64, graunlarity uint64, txn common.Transaction)

	SetupMergeOperators(db *kv.KeyValueDB, prefix string)

	Allocate(offset uint64, length uint64, txn common.Transaction)
}

type FreelistManager struct {
	freelistManagerInterface
	Cct *types.CephContext
}

func (fm *FreelistManager) New(cct *types.CephContext) {
	fm.Cct = cct
}

func CreateFreelistManage(cct *types.CephContext, _type string, kvdb *kv.KeyValueDB, prefix string) *FreelistManager {
	utils.AssertTrue("B" == prefix)
	if _type == "bitmap" {
		return &FreelistManager{
			freelistManagerInterface: NewBitmapFreelistManager(cct, kvdb, "B", "b"),
		}
	}
	return nil
}
