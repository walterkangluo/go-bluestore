package types

type MdConfigT struct {
	BlueFsAllocSize       uint64
	BlueFsSharedAllocSize uint64
	BlueFsAllocator       string
	BlueFsMaxLogRunaway   uint64

	BlueStoreBlueFs           bool
	BlueStoreBlueFsMinRation  float64
	BlueStoreBlueFsGiftRation float64
	BlueStoreBlueFsMin        uint64
	BlueStoreBlueFsEnvMirror  bool

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

	BlueStoreCacheSize              uint64
	BlueStoreCacheSizeHdd           uint64
	BlueStoreCacheSizeSSd           uint64
	BlueStoreCacheMetaRation        float64
	BlueStoreCacheKVRatio           float64
	BlueStoreCacheAutotune          bool
	BlueStoreCacheAutotuneChunkSize uint64
	BlueStoreCacheAutotuneInterval  float64
	BlueStoreOsdMemoryTarget        uint64
	OsdMemoryBase                   uint64
	OsdMemoryTarget                 uint64
	OsdMemoryExpectedFragmentation  float64
	OsdCacheCacheMin                uint64
	OsdMemoryCacheResizeInterval    uint64

	BlueStoreMinAllocSize    uint64
	BlueStoreMinAllocSizeHdd uint64
	BlueStoreMinAllocSizeSSd uint64

	BlueStoreDebugPermitAnyBdevLabel bool
	OsdMaxObjectSize                 uint32

	BlueStoreKVBackend string

	BdevAioMaxQueueDepth int
	Perf                 bool
	RockDBSeperateWalDir bool
}

func (md *MdConfigT) GetVal(key interface{}) interface{} {
	switch key.(type) {
	}
	return nil
}
