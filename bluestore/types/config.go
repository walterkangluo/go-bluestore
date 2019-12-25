package types

type MdConfigT struct {
	BlueFsAllocSize               uint64
	BlueFsSharedAllocSize         uint64
	BlueFsAllocator               string
	BlueFsMaxLogRunaway           uint64
	BlueStoreFsckOnMount          bool
	BlueStoreFsckOnMountDeep      bool
	BlueStoreFsckOnMkfs           bool
	BlueStoreFsckOnMkfsDeep       bool
	BlueStoreBlockPreallocateSize bool
	BlueStoreBlockPath            string
	BlueStoreBlockSize            uint64
	BlueStoreBlockCreate          bool
	BlueStoreBlockWalPath         string
	BlueStoreBlockWalSize         uint64
	BlueStoreBlockWalCreate       bool
	BlueStoreBlockDbPath          string
	BlueStoreBlockDbSize          uint64
	BlueStoreBlockDbCreate        bool
	BlueStoreBlueFs               bool
	OsdMaxObjectSize              uint32
}
