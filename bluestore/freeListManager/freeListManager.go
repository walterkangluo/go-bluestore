package freeListManager

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/utils"
)

type freeListManage interface {
}

type FreelistManager struct {
	Cct *types.CephContext
}

func CreateFreelistManager(cct *types.CephContext, _type string, kvdb types.KeyValueDB, prefix string) *FreelistManager {
	utils.AssertTrue("B" == prefix)
	if _type == "bitmap" {

	}
	return nil
}
