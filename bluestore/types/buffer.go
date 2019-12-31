package types

import (
	"fmt"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/utils"
	"os"
)

type Raw struct {
	data []byte //data stored
	off  uint64 //
	len  uint64 //存储的数据长度，0也算
}

type BufferList struct {
	ptr    []Raw
	size   uint64   //bufferlist size
	len    uint64   //数据总长度单位为byte
	delete []uint64 //已删除数据的下标列表，代表，buffer中的空洞
}

//Raw func
func (r *Raw) GetOff() uint64 {
	return r.off
}

func (r *Raw) Length() uint64 {
	return uint64(len(r.data))
}

func (bl *BufferList) Init() {
	bl.ptr = make([]Raw, 0)
	bl.size = 0
	bl.len = 0
	bl.delete = make([]uint64, 0)
}

func (bl *BufferList) Add(data []byte, dLen uint64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if len(bl.delete) != 0 {
		i := bl.delete[0]
		if bl.ptr[i].data != nil {
			panic("can not cover a exist Raw data")
		}
		bl.ptr[i].data = data
		bl.ptr[i].len = dLen
		bl.len += dLen

		bl.delete = bl.delete[1:]
	} else {
		raw := Raw{data, bl.size, dLen}
		bl.ptr = append(bl.ptr, raw)
		bl.size++
		bl.len += dLen
	}
}

func (bl *BufferList) AppendZero(len uint64) {
	zero := make([]byte, len)
	bl.Add(zero, len)
}

func (bl *BufferList) Delete(pos uint64) {
	if bl.ptr[pos].data == nil {
		return
	}

	bl.ptr[pos].data = nil
	bl.len -= bl.ptr[pos].len
	bl.delete = append(bl.delete, pos)
}

func (bl *BufferList) Begin() *Raw {
	return &bl.ptr[0]
}

func (bl *BufferList) End() *Raw {
	return &bl.ptr[len(bl.ptr)-1]
}

func (bl *BufferList) GetAt(pos uint64) *Raw {
	return &bl.ptr[pos]
}

func (bl *BufferList) Length() uint64 {
	return bl.size
}

func (bl *BufferList) SetLenth(len uint64) {

}

//buffer list other functions
func (bl *BufferList) ReadFd(fd *os.File, len uint64) int64 {
	var buf string
	ret, err := common.SafeRead(fd, &buf, int64(len))
	if err == nil {
		bl.Add([]byte(buf), uint64(ret))
	}
	return ret
}

func (bl *BufferList) Decode(data []byte, r *Raw) {
	r.data = data
}

func (bl *BufferList) Encode(data []byte) {
}

func (bl *BufferList) SubstrOf(other *BufferList, off uint64, len uint64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if off+len > other.Length() {
		panic("end of buffer")
	}

	var i uint64
	for i = 0; i < other.Length(); i++ {
		if off > 0 && off >= other.GetAt(i).Length() {
			off -= other.GetAt(i).Length()
		}
	}
	utils.AssertTrue(len == 0 || i != other.Length())

	for len > 0 {
		if off+len < other.ptr[i].Length() {
			bl.Add(other.ptr[i].data[off:], len)
			break
		}

		howMuch := other.ptr[i].Length() - off
		bl.len += howMuch
		len -= howMuch
		off = 0
		i++
	}
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
