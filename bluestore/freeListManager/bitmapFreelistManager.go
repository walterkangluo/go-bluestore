package freeListManager

import (
	"github.com/go-bluestore/bluestore/kv/keyvalue_db"
	"github.com/go-bluestore/bluestore/types"
	"sync"
)

type BitmapFreelistManager struct {
	FreelistManager
	kvdb *keyvalue_db.KeyValueDB
	Lock sync.Mutex

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
