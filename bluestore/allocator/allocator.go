package allocator

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
)

func CreateAllocator(cct *types.CephContext, _type string, size int64, blockSize int64, name string) Allocator {
	var alloc Allocator = nil

	if _type == "stupid" {
		alloc = CreateStupidAllocator(cct, name)
	} else if _type == "bitmap" {
		alloc = CreateBitmapAllocator(cct, size, blockSize, name)
	} else {
		log.Error("unknown type %s.", _type)
		panic("unknown type")
	}

	return alloc
}
