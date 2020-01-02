package freeListManager

import (
	"github.com/go-bluestore/bluestore/kv/keyvalue_db"
	"github.com/go-bluestore/bluestore/types"
	"sync"
)

type BitmapFreelistManager struct {
	*FreelistManager
	metaPrefix   string
	bitmapPrefix string
	kvdb         keyvalue_db.KeyValueDBI
	Lock         sync.Mutex

	size          uint64
	bytesPerBlock uint64
	blocksPerKey  uint64
	bytesPerKey   uint64
	blocks        uint64
	keyMask       uint64

	allSetBl        types.BufferList
	enumerateOffset uint64
	enumerateBlPos  int
}

func CreateBitmapFreelistManager(cct *types.CephContext, db keyvalue_db.KeyValueDBI,
	metaPrefix string, bitmapPrefix string) (bmfm *BitmapFreelistManager) {

	bmfm = &BitmapFreelistManager{
		metaPrefix:     metaPrefix,
		kvdb:           db,
		enumerateBlPos: 0,
		bitmapPrefix:   bitmapPrefix,
	}
	bmfm.FreelistManager.New(cct)

	return
}
