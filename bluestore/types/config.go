package types

type MdConfigT struct {
	BlueFsAllocSize          uint64
	BlueFsSharedAllocSize    uint64
	BlueFsAllocator          string
	BlueFsMaxLogRunaway      uint64
	BlueStoreFsckOnMount     bool
	BlueStoreFsckOnMountDeep bool
	OsdMaxObjectSize         uint32
}
