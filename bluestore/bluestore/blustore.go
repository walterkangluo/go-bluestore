package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/lib/thread_pool"
	"github.com/go-bluestore/log"
)

const (
	ObjectMaxSize = 0xffffffff
)

type SbInfoT struct {
	oidS       []types.GhObject
	sb         SharedBlob
	refMap     BlueStoreExtentRefMapT
	compressed bool
}

func (bs *BlueStore)ReadMeta(key string, value *string) int{
	return 0
}

func (bs *BlueStore)openPath() int {
	return 0
}

func (bs *BlueStore)openFsid(create bool) int {
	return 0
}

func (bs *BlueStore)readFsid(uuid types.UuidD) int {
	return 0
}

func (bs *BlueStore)lockFsid() int {
	return 0
}

func (bs *BlueStore)openBdev(create bool) int {
	return 0
}

func (bs *BlueStore)openDb(create bool) int {
	return 0
}

func (bs *BlueStore)openSuperMeta() int {
	return 0
}

func (bs *BlueStore)openFm(create bool) int {
	return 0
}

func (bs *BlueStore)openAlloc() int {
	return 0
}

func (bs *BlueStore)openCollections() int {
	return 0
}

func (bs *BlueStore)reloadLogger() int {
	return 0
}

func (bs *BlueStore)reconcileBluefsFreespace() int {
	return 0
}

func (bs *BlueStore)kvStart() int {
	return 0
}

func (bs *BlueStore)deferredReplay() int {
	return 0
}

func (bs *BlueStore)kvStop() int {
	return 0
}

func (bs *BlueStore)flushCache() int {
	return 0
}

func (bs *BlueStore)closeAlloc() int {
	return 0
}

func (bs *BlueStore)closeFm() int {
	return 0
}

func (bs *BlueStore)closeDb() int {
	return 0
}

func (bs *BlueStore)closeBdev() int {
	return 0
}

func (bs *BlueStore)closeFsid() int {
	return 0
}

func (bs *BlueStore)closePath() int {
	return 0
}

func (bs *BlueStore)fsck(deep bool, repair bool) int {
	return 0
}

func (bs *BlueStore)Fsck(deep bool)int{
	return bs.fsck(deep, false)
}

func (bs *BlueStore)mount(kvOnly bool) int{
	log.Debug("path %s", bs.Path)

	bs.KvOnly = kvOnly

	var mType string
	r := bs.ReadMeta("type", &mType)
	if r < 0 {
		log.Error("expected bluestore, but type is %s", mType)
		return -5
	}

	if mType != "bluestore" {
		log.Error("expected bluestore, but type is %s", mType)
		return -5
	}

	if bs.Cct.Conf.BlueStoreFsckOnMount{
		rc := bs.Fsck(bs.Cct.Conf.BlueStoreFsckOnMountDeep)
		if rc < 0{
			return rc
		}
		if rc > 0{
			log.Error("fsck found %d errors", rc)
			return -5
		}
	}

	if bs.Cct.Conf.OsdMaxObjectSize > ObjectMaxSize{
		log.Error("osd_max_object_size %u > bluestore max", bs.Cct.Conf.OsdMaxObjectSize)
		return -22
	}

	r = bs.openPath()
	if r < 0{
		return r
	}
	r = bs.openFsid(false)
	if r < 0{
		goto outPath
	}

	r = bs.readFsid(&bs.Fsid)
	if r < 0{
		goto outFsid
	}

	r = bs.openBdev(false)
	if r < 0{
		goto outFsid
	}

	r = bs.openDb(false)
	if r < 0{
		goto outBdev
	}

	if kvOnly{
		return 0
	}

	r = bs.openSuperMeta()
	if r < 0{
		goto outDb
	}

	r = bs.openFm(false)
	if r < 0{
		goto outDb
	}

	r = bs.openAlloc()
	if r < 0{
		goto outFm
	}

	r = bs.openCollections()
	if r < 0{
		goto outAlloc
	}

	r = bs. reloadLogger()
	if r < 0{
		goto outColl
	}

	if bs.BlueFS != nil{
		r = bs.reconcileBluefsFreespace()
		if r < 0{
			goto outColl
		}
	}

	bs.kvStart()

	r = bs.deferredReplay()
	if r < 0 {
		goto outStop
	}

	bs.MemPoolThread.New("bstore_mempool", 10, 0)

	bs.Mounted = true

	return 0

outStop:
	bs.kvStop()
outColl:
	bs.flushCache()
outAlloc:
	bs.closeAlloc()
outFm:
	bs.closeFm()
outDb:
	bs.closeDb()
outBdev:
	bs.closeBdev()
outFsid:
	bs.closeFsid()
outPath:
	bs.closePath()
	return r
}

func (bs *BlueStore)Mount() int{
	return bs.mount(false)
}