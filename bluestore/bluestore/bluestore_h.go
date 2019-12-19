package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/bluefs"
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

type TransContext struct {
	AioContext
}

type BlueStore struct {
	types.ObjectStore
	bluefs.BlueFSDeviceExpander
	types.MdConfigT

	ac AioContext
	tc TransContext
}

type Onode struct {
	onode         types.BlueStoreOnode
	exists        bool
	nref          int
	flushingCount int
	flushLock     sync.Mutex
	flushCond     *sync.Cond
	// onodeRef Onode
	// extentMap ExtentMap
}

func CreateOnode() *Onode {
	return &Onode{}
}

func (on *Onode) Flush() {

}

func (on *Onode) Get() {
	on.nref++
}

func (on *Onode) Put() {
	on.nref--
	if 0 == on.nref {
		on = new(Onode)
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