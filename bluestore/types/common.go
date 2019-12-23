package types

import (
	"github.com/go-bluestore/common/types"
	"github.com/go-bluestore/utils"
)

type BufferList struct {
	types.Vector
}

func (bl *BufferList) CreateVector(_preAlloc bool, _size int) {
	bl.CreateVector(_preAlloc, _size)
}

func (bl *BufferList) Init() {
	bl.Init()
}

func (bl *BufferList) Length() uint64 {
	return bl.Length()
}

func (bl *BufferList) Encode(data interface{}) []byte {
	return nil
}

func (bl *BufferList) AppendZero(length uint64) {
}

func (bl *BufferList) CRC32(src interface{}) uint32 {

	switch src.(type) {
	case int:
		if src.(int) == -1 {
			return utils.CRC32Byte(nil)
		}
	default:
		panic("not support")
	}
	return 0
}

func (bl *BufferList) Begin () {
	return bl.Begin()
}
