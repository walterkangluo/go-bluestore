package bluestore

import (
	"fmt"
	btypes "github.com/go-bluestore/bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/common"
	ctypes "github.com/go-bluestore/common/types"
	"github.com/go-bluestore/log"
	"os"
	"syscall"
	"unsafe"
	"github.com/go-bluestore/utils"
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

func (bs *BlueStore) ReadMeta(_key string, _value *string) int {
	var label *btypes.BluestoreBdevLabelT

	p := bs.Path + "/block"
	r := bs.readBdevLabel(bs.Cct, p, label)
	if r < 0 {
		return bs.ObjectStore.ReadMeta(_key, _value)
	}

	e, ok := label.Meta.Exists(ctypes.CreateElements(_key, nil))
	if !ok {
		return bs.ObjectStore.ReadMeta(_key, _value)
	}

	*_value = e.(*ctypes.Elements).GetVal().(string)
	return 0
}


func (bs *BlueStore) WriteMeta(_key string, _value string) int {
	var label *btypes.BluestoreBdevLabelT

	p := bs.Path + "/block"
	r := bs.readBdevLabel(bs.Cct, p, label)
	if r < 0 {
		return bs.ObjectStore.WriteMeta(_key, _value)
	}

	label.Meta.Push(ctypes.CreateElements(_key, nil))
	r = bs.writeBdevLabel(bs.Cct, p, label)
	utils.AssertTrue(r == 0)

	return bs.ObjectStore.WriteMeta(_key, _value)
}

func (bs *BlueStore) readBdevLabel(cct *types.CephContext, path string, label *btypes.BluestoreBdevLabelT) int {
	log.Debug("")

	var file *os.File
	for _, err := os.OpenFile(bs.Path, os.O_RDONLY|os.O_EXCL, 0); err != nil && err == syscall.EINTR; _, err = os.OpenFile(bs.Path, os.O_RDONLY|os.O_EXCL, 0) {
	}
	var bl types.BufferList
	r := bl.ReadFd(file, ObjectMaxSize)
	for err := file.Close(); err != nil && err == syscall.EINTR; err = file.Close() {
	}
	if r < 0 {
		log.Error("failed to read from %s: %d", path, r)
	}

	var crc, expectedCrc uint32
	p := bl.Front()
	defer func() {
		if err := recover(); err != nil {
			log.Debug("unable to decode label at offset %s")
			fmt.Println(err)
		}
	}()
	bl.Decode(*(*[]byte)(unsafe.Pointer(label)), p)
	var t types.BufferList
	/*TODO:暂时未想好实现，主要关于迭代器的问题,此处需要begin函数返回的对象为迭代器，现有实现中
	返回的是vector中的第一个元素，无法实现p.get_off()方法，待后续看这里的迭代器能否用其他用法代替*/
	//t.SubstrOf(bl, 0, p.get_off())
	crc = t.CRC32(-1)
	bl.Decode(*(*[]byte)(unsafe.Pointer(&expectedCrc)), p)

	if crc != expectedCrc {
		log.Error("bad crc on label, expected %d != actual %d", expectedCrc, crc)
	}
	log.Debug("got %v", *label)

	return 0
}

func (bs *BlueStore) writeBdevLabel(cct *types.CephContext, path string, label *btypes.BluestoreBdevLabelT) int {
	return 0
}

func (bs *BlueStore) openPath() int {
	return 0
}

func (bs *BlueStore) openFsid(create bool) int {
	return 0
}

func (bs *BlueStore) readFsid(uuid *types.UuidD) int {
	return 0
}

func (bs *BlueStore) lockFsid() int {
	return 0
}

func (bs *BlueStore) openBdev(create bool) int {
	return 0
}

func (bs *BlueStore) openDb(create bool) int {
	return 0
}

func (bs *BlueStore) openSuperMeta() int {
	return 0
}

func (bs *BlueStore) openFm(create bool) int {
	return 0
}

func (bs *BlueStore) openAlloc() int {
	return 0
}

func (bs *BlueStore) openCollections() int {
	return 0
}

func (bs *BlueStore) reloadLogger() int {
	return 0
}

func (bs *BlueStore) reconcileBluefsFreespace() int {
	return 0
}

func (bs *BlueStore) kvStart() int {
	return 0
}

func (bs *BlueStore) deferredReplay() int {
	return 0
}

func (bs *BlueStore) kvStop() int {
	return 0
}

func (bs *BlueStore) flushCache() int {
	return 0
}

func (bs *BlueStore) closeAlloc() int {
	return 0
}

func (bs *BlueStore) closeFm() int {
	return 0
}

func (bs *BlueStore) closeDb() int {
	return 0
}

func (bs *BlueStore) closeBdev() int {
	return 0
}

func (bs *BlueStore) closeFsid() int {
	return 0
}

func (bs *BlueStore) closePath() int {
	return 0
}

func (bs *BlueStore) fsck(deep bool, repair bool) int {
	return 0
}

func (bs *BlueStore) Fsck(deep bool) int {
	return bs.fsck(deep, false)
}

func (bs *BlueStore) mount(kvOnly bool) int {
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

	if bs.Cct.Conf.BlueStoreFsckOnMount {
		rc := bs.Fsck(bs.Cct.Conf.BlueStoreFsckOnMountDeep)
		if rc < 0 {
			return rc
		}
		if rc > 0 {
			log.Error("fsck found %d errors", rc)
			return -5
		}
	}

	if bs.Cct.Conf.OsdMaxObjectSize > ObjectMaxSize {
		log.Error("osd_max_object_size %d > bluestore max", bs.Cct.Conf.OsdMaxObjectSize)
		return -22
	}

	r = bs.openPath()
	if r < 0 {
		return r
	}
	r = bs.openFsid(false)
	if r < 0 {
		goto outPath
	}

	r = bs.readFsid(bs.Fsid)
	if r < 0 {
		goto outFsid
	}

	r = bs.openBdev(false)
	if r < 0 {
		goto outFsid
	}

	r = bs.openDb(false)
	if r < 0 {
		goto outBdev
	}

	if kvOnly {
		return 0
	}

	r = bs.openSuperMeta()
	if r < 0 {
		goto outDb
	}

	r = bs.openFm(false)
	if r < 0 {
		goto outDb
	}

	r = bs.openAlloc()
	if r < 0 {
		goto outFm
	}

	r = bs.openCollections()
	if r < 0 {
		goto outAlloc
	}

	r = bs.reloadLogger()
	if r < 0 {
		goto outColl
	}

	if bs.BlueFS != nil {
		r = bs.reconcileBluefsFreespace()
		if r < 0 {
			goto outColl
		}
	}

	bs.kvStart()

	r = bs.deferredReplay()
	if r < 0 {
		goto outStop
	}

	bs.MemPoolThread.New("bstore_mempool", 10, common.PoolFlags{})

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

func (bs *BlueStore) Mount() int {
	return bs.mount(false)
}


func (bs *BlueStore) MkFS() error {
	log.Debug("path is %s.", bs.Path)
	var r int
	if bs.Cct.Conf.OsdMaxObjectSize > ObjectMaxSize {
		log.Error("OsdMaxObjectSize %d size over ObjectMaxSize %d.", bs.Cct.Conf.OsdMaxObjectSize, ObjectMaxSize)
		return syscall.EINVAL
	}
	var done string
	r = bs.ReadMeta("mkfs_done", &done)
	if r == 0 {
		log.Debug("already make fs")
		if bs.Cct.Conf.BlueStoreFsckOnMkfs {
			r = bs.Fsck(bs.Cct.Conf.BlueStoreFsckOnMkfsDeep)
			if r < 0 {
				log.Error("fsck on mkfs found fatal error %d.", r)
				return syscall.Errno(r)
			}
			if r > 0 {
				log.Error("fsck found %d error.", r)
				return syscall.Errno(r)
			}
		}
		return syscall.Errno(r)
	}

	var btype string
	r = bs.ReadMeta("type", &btype)
	if r == 0 {
		if "bluestore" != btype {
			log.Error("expect type is bluestore, while type is %s.", btype)
			return syscall.EIO
		}
	} else {
		r = bs.WriteMeta("type", "bluestore")
		if r < 0 {
			return syscall.Errno(r)
		}
	}

	return syscall.Errno(0)
}