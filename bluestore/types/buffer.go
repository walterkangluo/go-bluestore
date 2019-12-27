package types

import (
	"fmt"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/utils"
	"os"
)

type Raw struct{
	data []byte
	off uint64
}

func (r *Raw)GetOff()uint64{
	return r.off
}

type BufferList struct {
	ptr []Raw
	size uint64	//bufferlist size
	delete []uint64
}

func (bl *BufferList)Init(){
	bl.ptr = make([]Raw, 0)
	bl.size = 0
	bl.delete = make([]uint64, 0)
}

func (bl *BufferList) Add(data []byte) {
	defer func(){
		if err := recover(); err != nil{
			fmt.Println(err)
		}
	}()
	if len(bl.delete) != 0{
		i := bl.delete[0]
		if bl.ptr[i].data != nil{
			panic("can not cover a exist Raw data")
		}
		bl.ptr[i].data = data

		bl.delete = bl.delete[1:]
	}else{
		raw := Raw{data, bl.size}
		bl.ptr = append(bl.ptr, raw)
		bl.size++
	}
}

func (bl *BufferList) AppendZero(len uint64){
	zero := make([]byte, len)
	bl.Add(zero)
}

func (bl *BufferList)Delete(pos uint64){
	if bl.ptr[pos].data == nil{
		return
	}

	bl.ptr[pos].data = nil
	bl.delete = append(bl.delete, pos)
}

func (bl *BufferList)Begin()*Raw{
	return &bl.ptr[0]
}

func (bl *BufferList)GetAt(pos uint64)*Raw{
	return &bl.ptr[pos]
}

func (bl *BufferList) Length()uint64{
	return bl.size
}

func (bl *BufferList)SetLenth(len uint64){

}

func (bl *BufferList) ReadFd(fd *os.File, len uint64) int64 {
	var buf string
	ret, err := common.SafeRead(fd, &buf, int64(len))
	if err == nil{
		bl.Add([]byte(buf))
	}
	return ret
}

func (bl *BufferList) Decode(data []byte, r *Raw){
	r.data = data
}

func (bl *BufferList)SubstrOf(other *BufferList, off uint64, len uint64){

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

