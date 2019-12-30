package types

import (
	"github.com/go-bluestore/bluestore/types"
	ctypes "github.com/go-bluestore/common/types"
	"github.com/go-bluestore/utils"
)

type BlueFsExtentT struct {
	Bedv   uint8
	Offset uint64
	Length uint32
}

func CreateBlueFsExtentT(b uint8, o uint64, l uint32) *BlueFsExtentT {
	return &BlueFsExtentT{
		Bedv:   b,
		Offset: o,
		Length: l,
	}
}

func (be *BlueFsExtentT) End() uint64 {
	return be.Offset + uint64(be.Length)
}

func (be BlueFsExtentT) Equal(b *BlueFsExtentT) bool {

	if be.Length == b.Length && be.Offset == b.Offset && be.Bedv == b.Bedv {
		return true
	}
	return false
}

type BlueFsFnodeT struct {
	Ino  uint64
	Size uint64
	//mtime time.Time
	PreferBdev uint8
	Extents    *ctypes.Vector // BlueFsExtentT
	Allocated  uint64
}

func CreateBlueFsFnodeT() *BlueFsFnodeT {
	bf := &BlueFsFnodeT{
		Ino:        uint64(0),
		Size:       uint64(0),
		PreferBdev: uint8(0),
		Allocated:  uint64(0),
	}
	bf.Extents.Init()
	return bf
}

func (bf *BlueFsFnodeT) GetAllocated() uint64 {
	return bf.Allocated
}

func (bf *BlueFsFnodeT) RecalculateAllocated() {
	bf.Allocated = uint64(0)

	for i := 0; i < bf.Extents.Size(); i++ {
		bf.Allocated += uint64(bf.Extents.At(i).(BlueFsExtentT).Length)
	}

}

func (bf *BlueFsFnodeT) AppendExtent(ext *BlueFsExtentT) {
	for i := 0; i < bf.Extents.Size(); i++ {
		if bf.Extents.At(i).(BlueFsExtentT).Equal(ext) {
			return
		}
	}

	bf.Allocated += uint64(ext.Length)
	bf.Extents.PushBack(ext)
}

// TODO: add other method

type BlueFsSuperT struct {
	Uuid      types.UUID
	OsdUuid   types.UUID
	Version   uint64
	BlockSize uint32
	LogFnode  BlueFsFnodeT
}

func CreateBlueFsSuperT() *BlueFsSuperT {
	return &BlueFsSuperT{
		Version:   uint64(0),
		BlockSize: uint32(4096),
	}
}

func (bs *BlueFsSuperT) Init() {
	bs.Version = uint64(0)
	bs.BlockSize = uint32(4096)
}

func (bs *BlueFsSuperT) blockMask() uint64 {
	return ^(uint64(bs.BlockSize) - uint64(1))
}

type opT uint

const (
	_ opT = iota
	opInit
	opAllocAdd
	opAllocRm
	opDirLink
	opDirUnlink
	opDirCreate
	opDirRemove
	opFileUpdate
	opFileRemove
	opJump
	opJumpSeq
)

type BlueFsTransactionT struct {
	Uuid types.UUID
	Seq  uint64
	opBl types.BufferList
}

func (bt *BlueFsTransactionT) Empty() bool {
	return bt.opBl.Length() == 0
}

func (bt *BlueFsTransactionT) OpInit() {
}

func (bt *BlueFsTransactionT) OpAllocAdd(id uint, start uint64, len uint64) {
}

func (bt *BlueFsTransactionT) OpDirCreate(dir string) {
	bt.opBl.Encode(utils.NumToBytes(opDirCreate))
	bt.opBl.Encode([]byte(dir))
}
