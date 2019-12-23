package types

import (
	"github.com/go-bluestore/utils"
)

type BufferList struct {
	data []byte
	size uint64
}

func CreateBufferList() *BufferList {
	return &BufferList{
		data: make([]byte, 0),
		size: 0,
	}
}

func (bf *BufferList) Init() {
	bf.data = make([]byte, 0)
	bf.size = 0
}

func (bf *BufferList) Length() uint64 {
	return bf.size
}

func (bf *BufferList) Encode(data []byte) []byte {
	for i := uint64(0); i < bf.size; i++ {
		bf.data = append(bf.data, data[i])
		bf.size++
	}
	return bf.data
}

func (bf *BufferList) AppendZero(length uint64) {
	for i := uint64(0); i < length; i++ {
		bf.data = append(bf.data, 0)
		bf.size++
	}
}

func (bf *BufferList) CRC32(src interface{}) uint32 {

	switch src.(type) {
	case int:
		if src.(int) == -1 {
			return utils.CRC32Byte(bf.data)
		}
	default:
		panic("not support")
	}
	return 0
}
