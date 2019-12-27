package bluefs

import (
	al "github.com/go-bluestore/bluestore/allocator"
	"github.com/go-bluestore/bluestore/blockdevice"
	btypes "github.com/go-bluestore/bluestore/bluefs/types"
	types2 "github.com/go-bluestore/bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"sync"
	"syscall"
)

func (bf *BlueFS) SetSlowDeviceExpander(bfe *BlueFSDeviceExpander) {
	bf.slowDevExpander = bfe
}

func (bf *BlueFS) AddBlockDevice(id int, path string) error {
	log.Debug("bdev id %d and path %s.", id, path)

	utils.AssertTrue(id < bf.bdev.Size())
	utils.AssertTrue(nil == bf.bdev.At(id))

	b := blockdevice.CreateBlockDevice(bf.Cct, path, nil, nil)

	r := b.Open(path)
	if nil != r {
		return r
	}

	log.Debug("bdev:%d, path:%s, size:%d.", id, path, b.GetSize())

	bf.bdev.SetAt(id, b)
	bf.ioc.SetAt(id, blockdevice.CreateIOContext(bf.Cct, nil, false))
	return nil
}

func (bf *BlueFS) BdevSupportLabel(id int) bool {
	utils.AssertTrue(id < bf.bdev.Size())
	utils.AssertTrue(nil != bf.bdev.At(int(id)))

	return bf.bdev.At(id).(*blockdevice.BlockDevice).SupportedBdevLable()
}

func (bf *BlueFS) AddBlockExtent(id int, offset uint64, length uint64) {
	log.Debug("bdev id %d offset %d length %d.", id, offset, length)

	utils.AssertTrue(id < bf.bdev.Size())
	utils.AssertTrue(nil != bf.bdev.At(id))
	utils.AssertTrue(bf.bdev.At(id).(*blockdevice.BlockDevice).GetSize() >= offset+length)

	bf.blockAll[id].insert(offset, length)
	bf.blockTotal[id] += length

	var l sync.Mutex
	if id < bf.alloc.Size() && nil != bf.alloc.At(id) {
		bf.logT.OpAllocAdd(uint(id), offset, length)
		r := bf.flushAndSyncLog(&l, 0, 0)
		utils.AssertTrue(r == nil)
		bf.alloc.At(id).(al.Allocator).InitAddFree(offset, length)
	}

	if nil != bf.logger {
		bf.logger.Inc(lBlueFsGiftBytes, length)
	}

	log.Debug("done")
}

func (bf *BlueFS) GetBlockDeviceSize(deviceId int) uint64 {
	utils.AssertTrue(deviceId < bf.bdev.Size())
	utils.AssertTrue(nil != bf.bdev.At(int(deviceId)))

	if deviceId < bf.bdev.Size() && nil != bf.bdev.At(int(deviceId)) {
		return bf.bdev.At(deviceId).(*blockdevice.BlockDevice).GetBlockSize()
	}

	return 0
}

func initAlloc(bfs *BlueFS) {
	bfs.alloc.ReSize(MaxBdev)          // make([]al.Allocator, MaxBdev)
	bfs.allocSize.ReSize(MaxBdev)      //  = make([]uint64, 0, MaxBdev)
	bfs.pendingRelease.ReSize(MaxBdev) // = make([]uint64, 0, MaxBdev)

	if nil != bfs.bdev.At(BdevWal) {
		bfs.allocSize.SetAt(BdevWal, bfs.Cct.Conf.BlueFsAllocSize) // = bfs.Cct.Conf.BlueFsAllocSize
	}

	if nil != bfs.bdev.At(BdevSlow) { // [BdevSlow]
		bfs.allocSize.SetAt(BdevDb, bfs.Cct.Conf.BlueFsAllocSize)
		bfs.allocSize.SetAt(BdevSlow, bfs.Cct.Conf.BlueFsSharedAllocSize)
		//bfs.allocSize[BdevDb] = bfs.Cct.Conf.BlueFsAllocSize
		//bfs.allocSize[BdevSlow] = bfs.Cct.Conf.BlueFsSharedAllocSize
	} else {
		//bfs.allocSize[BdevDb] = bfs.Cct.Conf.BlueFsAllocSize
		bfs.allocSize.SetAt(BdevDb, bfs.Cct.Conf.BlueFsAllocSize)
	}

	var blueFsFile = []string{"bluefs-wal", "bluefs-db", "bluefs-slow"}
	for id := 0; id < bfs.bdev.Size(); id++ {
		if nil == bfs.bdev.At(id) {
			continue
		}
		utils.AssertTrue(bfs.bdev.At(id).(*blockdevice.BlockDevice).GetSize() > 0)
		utils.AssertTrue(bfs.allocSize.At(id).(uint64) > 0)

		log.Debug("bdev name %s, allocSize %d, size %d.",
			blueFsFile[id], bfs.allocSize.At(id).(uint64), bfs.bdev.At(id).(*blockdevice.BlockDevice).GetSize())

		allocator := al.CreateAllocator(bfs.Cct, bfs.Cct.Conf.BlueFsAllocator,
			int64(bfs.bdev.At(id).(*blockdevice.BlockDevice).GetSize()), bfs.allocSize.At(id).(int64), blueFsFile[id])
		bfs.alloc.SetAt(id, allocator)
		blockAll := bfs.blockAll[id]

		for i := uint32(0); i < blockAll.size; i++ {
			segment := blockAll.segment[i]
			log.Debug("index %d and block start is %d and end is %d.", i, segment.GetStart(), segment.GetLen())
			bfs.alloc.At(id).(al.Allocator).InitAddFree(segment.GetStart(), segment.GetLen())
		}
	}
}

func initLogger(bfs *BlueFS) {

}

func (bfs *BlueFS) allocate(id uint8, l uint64, node *btypes.BlueFsFnodeT) error {
	log.Debug("len %d form device type %d", l, id)
	utils.AssertTrue(int(id) < bfs.alloc.Size())
	var extents types2.PExtentVector
	var allocLen int64
	extents.Init()
	if nil != bfs.alloc.At(int(id)) {
		hInt := uint64(0)
		if !node.Extents.Empty() && node.Extents.Back().(btypes.BlueFsExtentT).Bedv == id {
			hInt = node.Extents.Back().(*btypes.BlueFsExtentT).End()
		}
		extents.Reserve(4)
		allocLen = bfs.alloc.At(int(id)).(al.Allocator).Allocate(
			uint64(utils.RoundUpTo(int64(l), bfs.allocSize.At(int(id)).(int64))), bfs.allocSize.At(int(id)).(int64), int64(hInt), &extents)
	}

	if nil == bfs.alloc.At(int(id)) || allocLen < 0 || allocLen < int64(utils.RoundUpTo(int64(l), int64(bfs.allocSize.At(int(id)).(int64)))) {
		if allocLen > 0 {
			bfs.alloc.At(int(id)).(al.Allocator).Release(extents)
		}
		if id != BdevSlow {
			if nil != bfs.bdev.At(int(id)) {
				log.Error("failed to allocate %d on bdev %d, free %x, fallback to bdev %d.",
					l, id, bfs.alloc.At(int(id)), id+1)
			}
			return bfs.allocate(id+1, l, node)
		}

		if nil != bfs.bdev.At(int(id)) {
			log.Error("failed to allocate %x on bdev %d with free %d.",
				l, id, bfs.alloc.At(int(id)).(al.Allocator).GetFree())
		} else {
			log.Error("failed to allocate %x on bdev %d.", l, id)
		}

		if nil != bfs.bdev.At(int(id)) {
			bfs.alloc.At(int(id)).(al.Allocator).Dump()
		}

		return syscall.ENOSPC
	}

	for i := 0; i < extents.Size(); i++ {
		pe := extents.At(i).(types2.BlueStoreIntervalT)
		node.AppendExtent(&btypes.BlueFsExtentT{
			Bedv:   id,
			Offset: pe.Offset,
			Length: uint32(pe.Length),
		})
	}

	return nil
}

func (bfs *BlueFS) flushAndSyncLog(l *sync.Mutex, wantSeq uint64, jumpTo uint64) error {
	for {
		if bfs.logFlushing {
			log.Debug("want_seq %d, log is currently flushing, waiting", wantSeq)
			utils.AssertTrue(jumpTo == 0)
			bfs.logCond.Wait()
		}
	}
	/*
		if wantSeq != 0 && wantSeq < bfs.logSeqStable {
			log.Debug("want_seq %d <= logSeqStable %d, done", wantSeq, bfs.logSeqStable)
			utils.AssertTrue(jumpTo == 0)
			return nil
		}

		if bfs.logT.Empty() && bfs.dirtyFiles.
	*/
}

func (bfs *BlueFS) writeSuper() {
	var bl types.BufferList
	bl.Init()
	//var crc = bl.CRC32(-1)

	utils.AssertTrue(bl.Length() <= getSuperLength())
	bl.AppendZero(getSuperLength() - bl.Length())
	bfs.bdev.At(BdevDb).(*blockdevice.BlockDevice).Write(getSuperLength(), *bl, false)
}

func (bfs *BlueFS) flushBdev() {
	for i := 0; i < bfs.bdev.Size(); i++ {
		bfs.bdev.At(i).(*blockdevice.BlockDevice).Flush()
	}
}

func (bfs *BlueFS) closeWriter() {
	var h = bfs.logWriter
	log.Debug("write type is %d.", h.writerType)

	for i := 0; i < MaxBdev; i++ {
		if nil != bfs.bdev.At(i) {
			utils.AssertTrue(nil != h.iocv[i])

			h.iocv[i].AioWait()
			bfs.bdev.At(i).(*blockdevice.BlockDevice).QueueReapIoc()
		}
	}
}

func (bfs *BlueFS) stopAlloc() {
	for i := 0; i < bfs.alloc.Size(); i++ {
		if p := bfs.alloc.At(i); p != nil {
			p.(al.Allocator).Shutdown()
		}
	}
	bfs.alloc.Clear()
}

func (bfs *BlueFS) shutdownLogger() {
	bfs.Cct.GetPerfCountersCollection().Remove(bfs.logger)
	bfs.logger = nil
}

func (bfs *BlueFS) Mkfs(osdUuid types.UUID) {
	log.Debug("osd uuid is %v", osdUuid.UUID)
	var l sync.Mutex
	initAlloc(bfs)

	initLogger(bfs)

	super := btypes.BlueFsSuperT{
		Version:   uint64(1),
		BlockSize: uint32(bfs.bdev.At(BdevDb).(*blockdevice.BlockDevice).GetBlockSize()),
		OsdUuid:   osdUuid,
		Uuid:      types.GenerateRandomUuid(),
	}

	log.Debug("super uuid is %v", super.Uuid)

	var logFile fileRef
	logFile.fnode.Ino = uint64(1)
	logFile.fnode.PreferBdev = BdevWal

	r := bfs.allocate(logFile.fnode.PreferBdev, bfs.Cct.Conf.BlueFsMaxLogRunaway, &logFile.fnode)
	utils.AssertTrue(nil == r)

	bfs.logT.OpInit() // TODO
	for i := 0; i < MaxBdev; i++ {
		p := bfs.blockAll[i]

		if p.empty() {
			continue
		}

		for j := uint32(0); j < p.size; j++ {
			seg := p.segment[j]
			log.Debug("op alloc add start[%x] and length[%x].", seg.GetStart(), seg.GetLen())
			bfs.logT.OpAllocAdd(uint(i), seg.GetStart(), seg.GetLen())
		}
	}

	bfs.flushAndSyncLog(&l, 0, 0)

	super.LogFnode = logFile.fnode
	bfs.writeSuper()
	bfs.flushBdev()

	bfs.super = btypes.CreateBlueFsSuperT()
	bfs.closeWriter()
	bfs.logWriter = nil
	//bfs.blockAll
	bfs.blockTotal = make([]uint64, 0)
	bfs.stopAlloc()
	bfs.shutdownLogger()

	log.Debug("make bluefs success")
}

func (bfs *BlueFS) Mount() error {
	return nil
}

func (bfs *BlueFS) MkDir(dirName string) error {
	log.Debug("dir name %s.", dirName)

	_, ok := bfs.dirMap[dirName]
	if ok {
		// dirname has exist
		log.Error("dir name %s has exists.", dirName)
		return syscall.EEXIST
	}

	var ref dirRef
	ref.New()
	bfs.dirMap[dirName] = ref
	bfs.logT.OpDirCreate(dirName)
	return nil
}
