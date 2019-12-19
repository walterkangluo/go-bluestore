package types

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
	ino  uint64
	size uint64
	//mtime time.Time
	preferBdev uint8
	extents    []BlueFsExtentT
	allocated  uint64
}

func CreateBlueFsFnodeT() *BlueFsFnodeT {
	return &BlueFsFnodeT{
		ino:        uint64(0),
		size:       uint64(0),
		preferBdev: uint8(0),
		allocated:  uint64(0),
	}
}

func (bf *BlueFsFnodeT) getAllocated() uint64 {
	return bf.allocated
}

func (bf *BlueFsFnodeT) recalculateAllocated() {
	bf.allocated = uint64(0)

	for _, val := range bf.extents {
		bf.allocated += uint64(val.Length)
	}
}

func (bf *BlueFsFnodeT) appendExtent(ext *BlueFsExtentT) {
	var key int
	var val BlueFsExtentT

	for key, val = range bf.extents {
		if val.Equal(ext) {
			break
		}
	}

	bf.allocated += uint64(val.Length)
	bf.extents[key] = *new(BlueFsExtentT)
}

// TODO: add other method
