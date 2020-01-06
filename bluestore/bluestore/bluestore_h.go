package bluestore

import (
	"github.com/go-bluestore/bluestore/allocator"
	"github.com/go-bluestore/bluestore/blockdevice"
	"github.com/go-bluestore/bluestore/bluefs"
	fm "github.com/go-bluestore/bluestore/freeListManager"
	"github.com/go-bluestore/bluestore/kv"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/lib/thread_pool"
	"github.com/golang/protobuf/ptypes/timestamp"
	"sync"
)

const (
	BypassCleanCache           = 0x1
	StateFsAllocated           = 0
	StateFsStored              = 1
	StateFsCompressedOriginal  = 2
	StateFsCompressed          = 3
	StateFsCompressedAllocated = 4
	StateFsLast                = 5

	SpdkPrefix = "spdk:"
)

type BuffSpace struct {
}

type CollectionRef Collection

type recordT struct {
	length uint32
	refs   uint32
}

func CreateRecordT(l uint32, r uint32) *recordT {
	return &recordT{
		length: l,
		refs:   r,
	}
}

type mapT map[uint64]recordT

type BlueStoreExtentRefMapT struct {
	refMap mapT
}

func (be *BlueStoreExtentRefMapT) Empty() bool {
	if 0 == len(be.refMap) {
		return true
	}
	return false
}

type Cache struct {
	cct        *types.CephContext
	numExtents uint64
	numBlobs   uint64
}

func (ca *Cache) addBlob() {
	ca.numBlobs++
}

func (ca *Cache) rmBlob() {
	ca.numBlobs--
}

func (ca *Cache) addExtent() {
	ca.numExtents++
}

func (ca *Cache) rmExtent() {
	ca.numExtents--
}

type Collection struct {
	cache *Cache
}

type BlueStoreSharedBlobT struct {
	SBId   uint64
	RefMap BlueStoreExtentRefMapT
}

func CreateBlueStoreSharedBlobT(_sBId uint64) *BlueStoreSharedBlobT {
	return &BlueStoreSharedBlobT{
		SBId: _sBId,
	}
}

func (bs *BlueStoreSharedBlobT) Empty() bool {
	return bs.RefMap.Empty()
}

func (bs *BlueStoreSharedBlobT) GenerateTestInstance(blobs []*BlueStoreSharedBlobT) {
	return
}

type SharedBlob struct {
	nRef         int
	loaded       bool
	coll         CollectionRef
	sBidUnLoaded uint64
	Persistent   *BlueStoreSharedBlobT
	bc           BuffSpace
}

func CreateSharedBlob2(i uint64, _coll *Collection) *SharedBlob {
	return &SharedBlob{}
}

func CreateSharedBlob(_coll *Collection) *SharedBlob {
	if nil != _coll.cache {
		_coll.cache.addBlob()
	}
	return nil
}

func (sb *SharedBlob) getCache() *Cache {
	if nil == sb.coll.cache {
		return nil
	}
	return sb.coll.cache
}

func (sb *SharedBlob) getSBid() uint64 {
	if sb.loaded {
		return sb.Persistent.SBId
	}
	return sb.sBidUnLoaded
}

func (sb *SharedBlob) get() {
	sb.nRef++
}

func (sb *SharedBlob) put() {

}

func (sb *SharedBlob) getRef(offset uint64, length uint32) {

}

func (sb *SharedBlob) putRef(offset uint64, length uint32) {

}

func (sb *SharedBlob) finishWrite(seq uint64) {

}

func SharedBlobEqual(l *SharedBlob, r *SharedBlob) bool {
	return l.getSBid() == r.getSBid()
}

func (sb *SharedBlob) isLoaded() bool {
	return sb.loaded
}

type AioContext struct {
}

func (ac *AioContext) aioFinish(bs *BlueStore) {

}

type TransContext struct {
	AioContext
}

type Extents struct {
	Start  uint64
	Length uint64
}

type BlueStore struct {
	types.ObjectStore
	bluefs.BlueFSDeviceExpander
	types.MdConfigT

	KvOnly        bool
	Mounted       bool
	Path          string
	MemPoolThread *thread_pool.MempoolThread

	// private
	blueFs                         *bluefs.BlueFS
	blueFsSharedBdev               int
	blueFsSingleSharedDevice       bool
	blueFsLastBalance              timestamp.Timestamp
	nextDumpOnBlueFsBalanceFailure timestamp.Timestamp

	db           *kv.KeyValueDB
	bdev         *blockdevice.BlockDevice
	freelistType string
	fm           *fm.FreelistManager
	alloc        *allocator.Allocator
	fsId         types.UUID
	pathFd       int
	fsIdFd       int
	mounted      bool

	collLock sync.RWMutex

	blueFsExtents           []Extents
	blueFsExtentsReclaiming []Extents

	blockSize      uint64
	blockMask      uint64
	blockSizeOrder uint
	minAllocSize   uint64

	// cache trim control
	cacheSize                      uint64
	cacheMetaRation                float64
	cacheKVRatio                   float64
	cacheDataRatio                 float64
	cacheAutotune                  bool
	cacheAutotuneChunkSize         uint64
	cacheAutotuneInterval          float64
	osdMemoryTarget                uint64
	osdMemoryBase                  uint64
	osdMemoryExpectedFragmentation float64
	osdCacheCacheMin               uint64
	osdMemoryCacheResizeInterval   uint64
}

type ONode struct {
	oNode         types.BlueStoreOnode
	exists        bool
	nRef          int
	flushingCount int
	flushLock     sync.Mutex
	flushCond     *sync.Cond
	// onodeRef Onode
	// extentMap ExtentMap
}

func CreateONode() *ONode {
	return &ONode{}
}

func (on *ONode) Flush() {

}

func (on *ONode) Get() {
	on.nRef++
}

func (on *ONode) Put() {
	on.nRef--
	if 0 == on.nRef {
		on = new(ONode)
	}
}

type volatileStatFs struct {
	values [StateFsLast]int64
}

func CreateVolatileStateFs() *volatileStatFs {
	return &volatileStatFs{}
}

func (vs *volatileStatFs) allocated() int64 {
	return vs.values[StateFsAllocated]
}

func (vs *volatileStatFs) stored() int64 {
	return vs.values[StateFsStored]
}

func (vs *volatileStatFs) compressedOriginal() int64 {
	return vs.values[StateFsCompressedOriginal]
}

func (vs *volatileStatFs) compressed() int64 {
	return vs.values[StateFsCompressed]
}

func (vs *volatileStatFs) compressedAllocated() int64 {
	return vs.values[StateFsCompressedAllocated]
}

func (vs *volatileStatFs) isEmpty() bool {
	if vs.values[StateFsAllocated] == int64(0) &&
		vs.values[StateFsStored] == int64(0) &&
		vs.values[StateFsCompressedOriginal] == int64(0) &&
		vs.values[StateFsCompressed] == int64(0) &&
		vs.values[StateFsCompressedAllocated] == int64(0) {
		return true
	}
	return false
}
