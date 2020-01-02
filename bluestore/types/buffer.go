package types

import (
	"fmt"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/utils"
	"hash/crc32"
	"os"
	"sync"
	"sync/atomic"
)

var (
	BufferTrackCrc          bool   = GetEnvBool("CEPH_BUFFER_TRACK")
	BufferCachedCrc         uint64 = 0
	BufferCachedCrcAdjusted uint64 = 0
	BufferMissedCrc         uint64 = 0
)

func GetEnvBool(key string) bool {
	val := os.Getenv(key)
	if len(val) == 0 {
		return false
	}
	if val == "off" || val == "no" || val == "false" || val == "0" {
		return false
	}
	return true
}

type ofsT struct {
	begin uint64
	end   uint64
}

type ccrcT struct {
	base uint32
	crc  uint32
}

type Raw struct {
	data        []byte //data stored
	off         uint64 //
	len         uint64 //存储的数据长度，0也算
	CrcMap      map[ofsT]ccrcT
	crcSpinlock *sync.Mutex
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
	return r.len
}

func (r *Raw) GetCrc(fromTo ofsT, crc *ccrcT) bool {
	r.crcSpinlock.Lock()
	i, ok := r.CrcMap[fromTo]
	if ok {
		*crc = i
		r.crcSpinlock.Unlock()
		return true
	} else {
		r.crcSpinlock.Unlock()
		return false
	}
}

func (r *Raw) SetCrc(fromTo *ofsT, crc *ccrcT) {
	r.crcSpinlock.Lock()
	r.CrcMap[*fromTo] = *crc
	r.crcSpinlock.Unlock()
}

//BufferList func
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
		raw := new(Raw)
		raw.CrcMap = make(map[ofsT]ccrcT)
		bl.ptr = append(bl.ptr, *raw)
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
	return bl.len
}

func (bl *BufferList) Size() uint64 {
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
	for i = 0; i < other.Size(); i++ {
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

func (bl *BufferList) CRC32(crc uint32) uint32 {
	for i := uint64(0); i < bl.Size(); i++ {
		r := bl.ptr[i]
		if r.Length() != 0 {
			var ofs = ofsT{r.GetOff(), r.GetOff() + r.Length()}
			var ccrc ccrcT
			if r.GetCrc(ofs, &ccrc) {
				if ccrc.base == crc {
					crc = ccrc.crc
					if BufferTrackCrc {
						atomic.AddUint64(&BufferCachedCrc, 1)
					}
				} else {
					crc = crc32.Update(crc, crc32.IEEETable, r.data)
					if BufferTrackCrc {
						atomic.AddUint64(&BufferCachedCrcAdjusted, 1)
					}
				}
			} else {
				if BufferTrackCrc {
					atomic.AddUint64(&BufferMissedCrc, 1)
				}
				var base = crc
				crc = crc32.Update(crc, crc32.IEEETable, r.data)
				var ccrc ccrcT
				ccrc.base = base
				ccrc.crc = crc
				r.SetCrc(&ofs, &ccrc)
			}
		}
	}
	return crc
}
