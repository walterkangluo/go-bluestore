package types

import "github.com/go-bluestore/bluestore/types"

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

func (be *BlueFsExtentT) Equal(b *BlueFsExtentT) bool {

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
	Extents    []BlueFsExtentT
	Allocated  uint64
}

func CreateBlueFsFnodeT() *BlueFsFnodeT {
	return &BlueFsFnodeT{
		Ino:        uint64(0),
		Size:       uint64(0),
		PreferBdev: uint8(0),
		Allocated:  uint64(0),
	}
}

func (bf *BlueFsFnodeT) getAllocated() uint64 {
	return bf.Allocated
}

func (bf *BlueFsFnodeT) recalculateAllocated() {
	bf.Allocated = uint64(0)

	for _, val := range bf.Extents {
		bf.Allocated += uint64(val.Length)
	}
}

func (bf *BlueFsFnodeT) appendExtent(ext *BlueFsExtentT) {
	var key int
	var val BlueFsExtentT

	for key, val = range bf.Extents {
		if val.Equal(ext) {
			break
		}
	}

	bf.Allocated += uint64(val.Length)
	bf.Extents[key] = *new(BlueFsExtentT)
}

// TODO: add other method

type BlueFsSuperT struct {
	Uuid      types.UuidD
	OsdUuid   types.UuidD
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

func (bs *BlueFsSuperT) blockMask() uint64 {
	return ^(uint64(bs.BlockSize) - uint64(1))
}

type BlueFsTransactionT struct {
	Uuid types.UuidD
	Seq  uint64
	opBl types.BufferList
}

func (bt *BlueFsTransactionT) Empty() bool {
	return bt.opBl.Length() == 0
}
