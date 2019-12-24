package types

import (
	"github.com/go-bluestore/common/types"
	"github.com/go-bluestore/utils"
	"os"
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

func (bl *BufferList) Decode(in []byte, data types.T) {

}

func (bl *BufferList) ReadFd(fd *os.File, len uint64) int64 {
	return 0
}

func (bl *BufferList) SubstrOf(other BufferList, off uint32, len uint32) {

}
