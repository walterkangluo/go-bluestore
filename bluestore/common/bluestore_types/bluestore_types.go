package bluestore_types

import (
	"github.com/go-bluestore/bluestore/common/formatter"
)

const FlagOmp = uint8(1)

type shardInfo struct {
	offset uint32
	bytes  uint32
}

type BlueStoreOnode struct {
	nid             uint64
	size            uint64
	attrs           map[string]interface{}
	extentMapShards []shardInfo

	expectedObjectSize uint32
	expectedWriteSize  uint32
	allocHintFlags     uint32
	flags              uint8
}

func (bo *BlueStoreOnode) getFlagsString() string {
	var s string
	m := bo.flags & FlagOmp
	if m != 0 {
		s = "omap"
	}
	return s
}

func (bo *BlueStoreOnode) hasFlags(f uint8) uint8 {
	return bo.flags & f
}

func (bo *BlueStoreOnode) setFlags(f uint8) {
	bo.flags |= f
}

func (bo *BlueStoreOnode) clearFlags(f uint8) {
	bo.flags &= f
}

func (bo *BlueStoreOnode) hasOmap() bool {
	m := bo.hasFlags(FlagOmp)

	if m != 0 {
		return true
	}

	return false
}

func (bo *BlueStoreOnode) setOmapFlag() {
	bo.setFlags(FlagOmp)
}

func (bo *BlueStoreOnode) clearOmapFlag() {
	bo.clearFlags(FlagOmp)
}

func (bo *BlueStoreOnode) dump(f formatter.Formatter) {

}

func (bo *BlueStoreOnode) generateTestInstance(o []*BlueStoreOnode) {

}
