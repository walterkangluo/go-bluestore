package bluestore

import (
	"fmt"
	"github.com/go-bluestore/bluestore/blockdevice"
	"github.com/go-bluestore/bluestore/bluefs"
	btypes "github.com/go-bluestore/bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/kv/keyvalue_db"
	"github.com/go-bluestore/bluestore/kv/rocksdb_store"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/common"
	ctypes "github.com/go-bluestore/common/types"
	lrdb "github.com/go-bluestore/lib/gorocksdb"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	PrefixSuper      = "S"
	PrefixState      = "T"
	PrefixColl       = "C"
	PrefixObj        = "O"
	PrefixOmap       = "M"
	PrefixDeferred   = "L"
	PrefixSAlloc     = "B"
	PrefixSharedBlob = "X"

	ObjectMaxSize      = 0xffffffff
	SuperReserved      = 8192
	BdevLabelBlockSize = 4096
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
	p := bl.Begin()
	defer func() {
		if err := recover(); err != nil {
			log.Debug("unable to decode label at offset")
			fmt.Println(err)
		}
	}()
	bl.Decode(*(*[]byte)(unsafe.Pointer(label)), p)
	var t types.BufferList
	t.SubstrOf(&bl, 0, p.GetOff())
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

func (bs *BlueStore) readFsid(uuid types.UUID) int {
	return 0
}

func (bs *BlueStore) lockFsid() int {
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
	var e error
	var mType string
	r := bs.ReadMeta("type", &mType)
	if r < 0 {
		log.Error("expected bluestore, but type is %s", mType)
		return r
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

	r = bs.readFsid(bs.fsId)
	if r < 0 {
		goto outFsid
	}

	e = bs.openBdev(false)
	if nil != e {
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

	if bs.blueFs != nil {
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

func (bs *BlueStore) setupBlockSymlinkOrFile(name string, epath string, size uint64, create bool) error {
	log.Debug("name: %s, path %s, size %d, create %v.", name, epath, size, create)

	var r error
	var flags = syscall.O_RDWR | syscall.O_CLOEXEC

	if create {
		flags |= syscall.O_CREAT
	}

	if 0 != len(epath) {
		r = syscall.Symlink(epath, name)
		if nil != r {
			log.Error("failed to create link for %s and %s.", epath, name)
			return r
		}

		if strings.HasPrefix(epath, SpdkPrefix) {
			file, r := os.OpenFile(epath, flags, 0644)
			if nil != r {
				log.Error("failed to open %s", epath)
				return r
			}
			defer file.Close()

			i := utils.Substr(epath, SpdkPrefix)
			utils.AssertTrue(i != -1)
			remainString := epath[i+len(SpdkPrefix):]
			n, r := file.WriteString(remainString)
			utils.AssertTrue(n == len(remainString))
			log.Debug("create %s symlink to %s.", name, epath)
			return r
		}
	}

	if size > 0 {
		file, r := os.OpenFile(epath, flags, 0644)
		if nil == r {
			st, r := file.Stat()
			if nil == r && st.Mode().IsRegular() && st.Size() == int64(0) {

				r = file.Truncate(int64(size))
				if nil != r {
					log.Error("failed to resize %s to %d.", name, size)
					return r
				}
			}

			if bs.Cct.Conf.BlueStoreBlockPreallocateSize {
				// TODO： implement fallcate manual
			}
			log.Debug("resize file %s to %d.", name, size)
		} else {
			log.Error("failed to open file %s.", name)
			return r
		}
	}
	return nil
}

func (bs *BlueStore) checkOrSetBdevLabel(path string, size uint64, desc string, create bool) error {
	var label btypes.BluestoreBdevLabelT
	var r int
	if create {
		label.OsdUUID = bs.fsId
		label.Size = size
		label.BTime = time.Now()
		label.Description = desc

		r = bs.readBdevLabel(bs.Cct, bs.Path, &label)
		if r < 0 {
			return fmt.Errorf("%d", r)
		}
	} else {
		r = bs.readBdevLabel(bs.Cct, bs.Path, &label)
		if r < 0 {
			return fmt.Errorf("%d", r)
		}

		if bs.Cct.Conf.BlueStoreDebugPermitAnyBdevLabel {
			log.Debug("bdev %s osdid %v fsid %v check passed.", bs.Path, label.OsdUUID, bs.fsId)
		} else {
			log.Error("bdev %s osdid %v does not match out fsid %v.", path, label.OsdUUID, bs.fsId)
			return syscall.EIO
		}
	}
	return nil
}

func (bs *BlueStore) setCacheSize() error {
	utils.AssertTrue(bs.bdev != nil)

	bs.cacheAutotune = bs.Cct.Conf.BlueStoreCacheAutotune
	bs.cacheAutotuneChunkSize = bs.Cct.Conf.BlueStoreCacheAutotuneChunkSize
	bs.cacheAutotuneInterval = bs.Cct.Conf.BlueStoreCacheAutotuneInterval
	bs.osdMemoryTarget = bs.Cct.Conf.OsdMemoryTarget
	bs.osdMemoryBase = bs.Cct.Conf.OsdMemoryBase
	bs.osdMemoryExpectedFragmentation = bs.Cct.Conf.OsdMemoryExpectedFragmentation
	bs.osdCacheCacheMin = bs.Cct.Conf.OsdCacheCacheMin
	bs.osdMemoryCacheResizeInterval = bs.Cct.Conf.OsdMemoryCacheResizeInterval

	if bs.Cct.Conf.BlueStoreCacheSize > 0 {
		bs.cacheSize = bs.Cct.Conf.BlueStoreCacheSize
	} else {
		if bs.bdev.BlockDeviceFunc.SupportedBdevLable() {
			bs.cacheSize = bs.Cct.Conf.BlueStoreCacheSizeHdd
		} else {
			bs.cacheSize = bs.Cct.Conf.BlueStoreCacheSizeSSd
		}
	}

	bs.cacheMetaRation = bs.Cct.Conf.BlueStoreCacheMetaRation
	if bs.cacheMetaRation < 0 || bs.cacheMetaRation > 1.0 {
		log.Error("BlueStoreCacheMetaRation must in range [0, 1.0]")
		return syscall.EINVAL
	}

	bs.cacheKVRatio = bs.Cct.Conf.BlueStoreCacheKVRatio
	if bs.cacheKVRatio < 0 || bs.cacheKVRatio > 1.0 {
		log.Error("BlueStoreCacheKVRatio must in range [0, 1.0]")
		return syscall.EINVAL
	}

	if bs.cacheMetaRation+bs.cacheKVRatio > 1.0 {
		log.Error("sum of BlueStoreCacheMetaRation and BlueStoreCacheKVRatio must in range [0, 1.0]")
		return syscall.EINVAL
	}

	bs.cacheDataRatio = 1.0 - bs.cacheMetaRation - bs.cacheKVRatio
	if 0 > bs.cacheDataRatio {
		bs.cacheDataRatio = 0
	}

	log.Debug("cache_size %d, meta %f, kv %f, data %f.",
		bs.cacheSize, bs.cacheMetaRation, bs.BlueStoreCacheKVRatio, bs.cacheDataRatio)

	return nil
}

func aioCb(priv unsafe.Pointer, priv2 unsafe.Pointer) {
	store := (*BlueStore)(priv)
	c := (*AioContext)(priv2)
	c.aioFinish(store)
}

func (bs *BlueStore) openBdev(create bool) error {
	var r error
	utils.AssertTrue(nil == bs.bdev)

	p := bs.Path + "/block"
	// TODO: implement create BlocDevice
	bs.bdev = blockdevice.CreateBlockDevice(bs.Cct, p, aioCb, unsafe.Pointer(bs))

	r = bs.bdev.BlockDeviceFunc.Open(p)
	if r != nil {
		log.Error("open path %s failed with %v.", p, r)
		goto fail
	}

	if bs.bdev.BlockDeviceFunc.SupportedBdevLable() {
		r = bs.checkOrSetBdevLabel(p, bs.bdev.BlockDeviceFunc.GetSize(), "main", create)
		if nil != r {
			goto failclose
		}
	}

	bs.blockSize = bs.bdev.BlockDeviceFunc.GetBlockSize()
	bs.blockMask = ^(bs.blockSize - 1)
	bs.blockSizeOrder = utils.Ctx(bs.blockMask)
	utils.AssertTrue(bs.blockSize == 1<<bs.blockSizeOrder)

	r = bs.setCacheSize()
	if r != nil {
		goto failclose
	}

	return nil

failclose:
	bs.bdev.BlockDeviceFunc.Close()
fail:
	bs.bdev = nil
	return r
}

func (bs *BlueStore) openDB(create bool) error {
	utils.AssertTrue(bs.db == nil)

	var r error
	var fn string
	var res string
	var options string
	var env *rocksdb_store.BlueRocksEnv

	// 1. get kv_backend type
	var kvBackend string
	if create {
		kvBackend = bs.Cct.Conf.BlueStoreKVBackend
	} else {
		r := bs.ReadMeta("kv_backend", &kvBackend)
		if r < 0 {
			log.Error("unable to read kv_backend meta.")
			return syscall.EIO
		}
	}

	log.Info("kv_backend = %s.", kvBackend)

	// 2. prepare bluefs for rockdb
	var doBlueFs bool
	if create {
		doBlueFs = bs.Cct.Conf.BlueStoreBlueFs
	} else {
		var blueFsFlag string
		r := bs.ReadMeta("bluefs", &blueFsFlag)
		if r < 0 {
			log.Error("unable to read bluefs meta.")
			return syscall.EIO
		}

		if "1" == blueFsFlag {
			doBlueFs = true
		} else if "0" == blueFsFlag {
			doBlueFs = false
		} else {
			log.Error("bluefs = %v, not 0 or 1, so aborting.", doBlueFs)
			return syscall.EIO
		}
	}

	log.Info("bluefs = %v.", doBlueFs)

	// prepare bluefs Env for rockdb
	if doBlueFs {
		log.Debug("initialing bluefs")
		if "rockdb" != kvBackend {
			log.Error("backend must be rockdb to use bluefs.")
			return syscall.EINVAL
		}

		bs.blueFs = bluefs.CreateBlueFS(bs.Cct)
		bs.blueFs.SetSlowDeviceExpander(&bs.BlueFSDeviceExpander)

		var bfn string
		var st syscall.Stat_t

		// block.db store meta data for bluestore
		bfn = bs.Path + "/block.db"
		if nil == syscall.Stat(bfn, &st) {
			r = bs.blueFs.AddBlockDevice(bluefs.BdevDb, bfn)
			if nil != r {
				log.Error("add block device %s failed with err %v.", bfn, r)
				goto freeBlueFs
			}

			if bs.blueFs.BdevSupportLabel(bluefs.BdevDb) {
				r = bs.checkOrSetBdevLabel(bfn, bs.blueFs.GetBlockDeviceSize(bluefs.BdevDb), "bluefs db", create)
				if nil != r {
					log.Error("check block device %s lable return %v.", bfn, r)
					goto freeBlueFs
				}
			}

			if create {
				bs.blueFs.AddBlockExtent(bluefs.BdevDb, SuperReserved, bs.blueFs.GetBlockDeviceSize(bluefs.BdevDb)-SuperReserved)
			}

			bs.blueFsSharedBdev = bluefs.BdevSlow
			bs.blueFsSingleSharedDevice = false
		} else {

			r = syscall.Lstat(bfn, &st)
			if r != nil {
				bs.blueFsSharedBdev = bluefs.BdevDb
			} else {
				log.Error("%s symlink exists but target unusable: %v.", bfn, r)
				goto freeBlueFs
			}
		}

		// block to store object data
		bfn = bs.Path + "/block"
		r = bs.blueFs.AddBlockDevice(bs.blueFsSharedBdev, bfn)
		if nil != r {
			log.Error("add block device %s return %v.", bfn, r)
			goto freeBlueFs
		}
		if create {
			initial := (float64(bs.bdev.BlockDeviceFunc.GetSize())) * (bs.Cct.Conf.BlueStoreBlueFsMinRation + bs.Cct.Conf.BlueStoreBlueFsGiftRation)
			initial = math.Max(initial, float64(bs.Cct.Conf.BlueStoreBlueFsMin))
			if 0 != bs.Cct.Conf.BlueFsAllocSize%bs.minAllocSize {
				r = syscall.EINVAL
				log.Error("bluefs_alloc_size %x is not mutiple of min_alloc_size %x.", bs.Cct.Conf.BlueFsAllocSize, bs.minAllocSize)
				goto freeBlueFs
			}

			initial = float64(utils.P2RoundUp(uint64(initial), bs.Cct.Conf.BlueFsAllocSize))
			start := utils.P2Align((bs.bdev.BlockDeviceFunc.GetSize()-uint64(initial))/2, bs.Cct.Conf.BlueFsAllocSize)
			bs.blueFs.AddBlockExtent(bs.blueFsSharedBdev, start, uint64(initial))
			bs.blueFsExtents = append(bs.blueFsExtents, Extents{start: start, length: uint64(initial)})
		}

		// block.wal to store log file of rockdb
		bfn = bs.Path + "/block.wal"
		r = syscall.Stat(bfn, &st)
		if r == nil {
			r = bs.blueFs.AddBlockDevice(bluefs.BdevWal, bfn)
			if r != nil {
				log.Error("add block device %s return %v.", bfn, r)
				goto freeBlueFs
			}

			if bs.blueFs.BdevSupportLabel(bluefs.BdevWal) {
				r = bs.checkOrSetBdevLabel(bfn, bs.blueFs.GetBlockDeviceSize(bluefs.BdevWal), "bluefs wal", create)
				if r != nil {
					log.Error("check block device %s label returned: %v.", bfn, r)
					goto freeBlueFs
				}
			}

			if create {
				bs.blueFs.AddBlockExtent(bs.blueFsSharedBdev, BdevLabelBlockSize, bs.blueFs.GetBlockDeviceSize(bluefs.BdevWal)-BdevLabelBlockSize)
			}

			bs.Cct.Conf.RocksDBSeperateWalDir = true
			bs.blueFsSingleSharedDevice = false
		} else {
			r = syscall.Lstat(bfn, &st)
			if r == nil {
				bs.Cct.Conf.RocksDBSeperateWalDir = false
			} else {
				log.Error("%s symlink exists, but target unusable: %v.", bfn, r)
				goto freeBlueFs
			}
		}

		if create {
			bs.blueFs.Mkfs(bs.fsId)
		}

		// TODO: implement mount
		r = bs.blueFs.Mount()
		if nil != r {
			log.Error("failed to mount: %v.", r)
			goto freeBlueFs
		}

		//var env *lrdb.Env
		if bs.Cct.Conf.BlueStoreBlueFsEnvMirror {
			a := rocksdb_store.NewBlueRocksEnv(bs.blueFs)
			b := lrdb.NewDefaultEnv()
			if create {
				cmd := "rm -rf " + bs.Path + "/db " + bs.Path + "/db.slow" + bs.Path + "/db.wal"
				res := exec.Command("sh", "-c", cmd)
				_, r := res.Output()
				utils.AssertTrue(r == nil)
			}
			env = rocksdb_store.NewEnvMirror(b, a.Wrapper, false, true)
		} else {
			env = rocksdb_store.NewBlueRocksEnv(bs.blueFs)
			fn = "db"
		}

		if bs.blueFsSharedBdev == bluefs.BdevSlow {
			// use block.db and block both to bluefs
			dbSize := bs.blueFs.GetBlockDeviceSize(bluefs.BdevDb)
			slowSize := bs.blueFs.GetBlockDeviceSize(bluefs.BdevSlow)
			dbPath := fmt.Sprintf("%s,%d %s.slow,%d", fn, uint64(float32(dbSize)*0.95), fn, uint64(float32(slowSize)*0.95))
			bs.Cct.Conf.RocksDBPaths = dbPath
			log.Debug("set rockdb_db_path to %s.", dbPath)
		}

		if create {
			log.Debug("%v and create dir %s.", env, fn)
			r = env.CreateDir(fn)
			utils.AssertTrue(r == nil)
			if bs.Cct.Conf.RocksDBSeperateWalDir {
				r = env.CreateDir(fn + ".wal")
				utils.AssertTrue(r == nil)
			}

			if len(bs.Cct.Conf.RocksDBPaths) != 0 {
				r = env.CreateDir(fn + ".slow")
				utils.AssertTrue(r == nil)
			}
		}
	} else {
		// do not use bluefs
		r = syscall.Mkdir(fn, 0755)
		if r != nil && r != syscall.EEXIST {
			log.Error("failed to create %s, with error %v.", fn, r)
			return r
		}

		if bs.Cct.Conf.RocksDBSeperateWalDir {
			walPath := bs.Path + "/db.wal"
			r = syscall.Mkdir(walPath, 0755)
			if r != nil && r != syscall.EEXIST {
				log.Error("failed to create %s, with error %v.", fn, r)
				return r
			}
		}
	}

	bs.db = keyvalue_db.CreateKeyValueDB(bs.Cct, kvBackend, fn, env)
	if nil == bs.db {
		log.Error("error create db")
		if nil != bs.blueFs {
			bs.blueFs.Umount()
			bs.blueFs = nil
			return syscall.EIO
		}
	}

	bs.db.SetMergeOperator(PrefixState)
	bs.db.SetCacheSize(uint64(bs.cacheKVRatio) * bs.cacheSize)

	if kvBackend == "rocksdb" {
		options = bs.Cct.Conf.BlueStoreRocksDBOptions
	}

	r = bs.db.Init(options)
	utils.AssertTrue(r == nil)

	if create {
		r = bs.db.CreateAndOpen(res)
	} else {
		r = bs.db.Open(res)
	}

	if nil != r {
		log.Error("error opening db")
		if nil != bs.blueFs {
			bs.blueFs.Umount()
			bs.blueFs = nil
			return syscall.EIO
		}
		bs.blueFs = nil
		return syscall.EIO
	}

	log.Error("opened %s path %s options %s.", kvBackend, fn, options)
	return nil

freeBlueFs:
	utils.AssertTrue(nil != bs.blueFs)
	bs.blueFs = nil

	return nil
}

func (bs *BlueStore) Mkfs() error {
	var r int
	var e error
	var oldFsId types.UUID

	log.Debug("path is %s.", bs.Path)

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

	r = bs.openPath()
	if r < 0 {
		return syscall.Errno(r)
	}

	r = bs.openFsid(true)
	if r < 0 {
		goto outPathFd
	}

	r = bs.lockFsid()
	if r < 0 {
		goto outCloseFsId
	}

	r = bs.readFsid(oldFsId)
	if r < 0 || oldFsId.IsZero() {
		if bs.fsId.IsZero() {
			bs.fsId = types.GenerateRandomUuid()
			log.Debug("generate fsid is %x.", bs.fsId)
		} else {
			log.Debug("using provided fsid %x.", bs.fsId)
		}
	} else {
		if !bs.fsId.IsZero() && bs.fsId != oldFsId {
			log.Error("ondisk uuid %x != provided %x.", oldFsId, bs.fsId)
			//r = -syscall.EINVAL
			r = -0x16
		}
		bs.fsId = oldFsId
	}

	e = bs.setupBlockSymlinkOrFile(
		"block", bs.Cct.Conf.BlueStoreBlockPath, bs.Cct.Conf.BlueStoreBlockSize, bs.Cct.Conf.BlueStoreBlockCreate)
	if nil != e {
		goto outCloseFsId
	}

	if bs.Cct.Conf.BlueStoreBlueFs {
		e = bs.setupBlockSymlinkOrFile(
			"block.wal", bs.Cct.Conf.BlueStoreBlockWalPath, bs.Cct.Conf.BlueStoreBlockWalSize, bs.Cct.Conf.BlueStoreBlockWalCreate)
		if nil != e {
			goto outCloseFsId
		}

		e = bs.setupBlockSymlinkOrFile(
			"block.db", bs.Cct.Conf.BlueStoreBlockDbPath, bs.Cct.Conf.BlueStoreBlockDbSize, bs.Cct.Conf.BlueStoreBlockDbCreate)
		if nil != e {
			goto outCloseFsId
		}
	}

	e = bs.openBdev(true)
	if e != nil {
		goto outCloseFsId
	}

	if bs.Cct.Conf.BlueStoreMinAllocSize > 0 {
		bs.minAllocSize = bs.Cct.Conf.BlueStoreMinAllocSize
	} else {
		utils.AssertTrue(nil != bs.bdev)
		if bs.bdev.BlockDeviceFunc.IsRotational() {
			bs.minAllocSize = bs.Cct.Conf.BlueStoreMinAllocSizeHdd
		} else {
			bs.minAllocSize = bs.Cct.Conf.BlueStoreMinAllocSizeSSd
		}
	}

	if !utils.ISP2(bs.minAllocSize) {
		log.Error("min_alloc_size %x is not power of 2 aligned!", bs.minAllocSize)
		e = syscall.EINVAL
		goto outCloseBdev
	}

outCloseBdev:
	bs.closeBdev()
outCloseFsId:
	bs.closeFsid()
outPathFd:
	bs.closePath()

	if r == 0 && e == nil && bs.Cct.Conf.BlueStoreFsckOnMkfs {
		rc := bs.Fsck(bs.Cct.Conf.BlueStoreFsckOnMkfsDeep)
		if rc < 0 {
			return syscall.Errno(r)
		}
		if rc > 0 {
			log.Error("found %d errors.", rc)
			e = syscall.EIO
		}
	}

	if r == 0 && e == nil {
		r = bs.WriteMeta("mkfs_done", "yes")
	}

	if r < 0 {
		log.Error("write mkfs_done failed with %d.", r)
	} else {
		log.Info("bluestore mkfs success")
	}

	return e
}
