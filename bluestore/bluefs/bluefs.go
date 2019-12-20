package bluefs

import (
	al "github.com/go-bluestore/bluestore/allocator"
	btypes "github.com/go-bluestore/bluestore/bluefs/types"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
)

func initAlloc(bfs *BlueFS) {
	bfs.alloc = make([]al.Allocator, MaxBdev)
	bfs.allocSize = make([]uint64, 0, MaxBdev)
	bfs.pendingRelease = make([]uint64, 0, MaxBdev)

	if nil != bfs.bdev[BdevWal] {
		bfs.allocSize[BdevWal] = bfs.Cct.Conf.BlueFsAllocSize
	}

	if nil != bfs.bdev[BdevSlow] {
		bfs.allocSize[BdevDb] = bfs.Cct.Conf.BlueFsAllocSize
		bfs.allocSize[BdevSlow] = bfs.Cct.Conf.BlueFsSharedAllocSize
	} else {
		bfs.allocSize[BdevDb] = bfs.Cct.Conf.BlueFsAllocSize
	}

	var blueFsFile = []string{"bluefs-wal", "bluefs-db", "bluefs-slow"}
	for id := 0; id < len(bfs.bdev); id++ {
		if nil == bfs.bdev[id] {
			continue
		}
		utils.AssertTrue(bfs.bdev[id].GetSize() > 0)
		utils.AssertTrue(bfs.allocSize[id] > 0)

		log.Debug("bdev name %s, allocSize %d, size %d.",
			blueFsFile[id], bfs.allocSize[id], bfs.bdev[id].GetSize())

		bfs.alloc[id] = al.CreateAllocator(
			bfs.Cct, bfs.Cct.Conf.BlueFsAllocator, int64(bfs.bdev[id].GetSize()), int64(bfs.allocSize[id]), blueFsFile[id])
		blockAll := bfs.blockAll[id]
		for index, blockInfo := range blockAll {
			log.Debug("index %d and block start is %d and end is %d.", index, blockInfo.getStart(), blockInfo.getLen())
			bfs.alloc[id].InitAddFree(blockInfo.getStart(), blockInfo.getLen())
		}
	}
}

func initLogger(bfs *BlueFS) {

}

func (bfs *BlueFS) allocate(id uint8, len uint64, node *btypes.BlueFsFnodeT) int {
	log.Debug("len %d form device type %d", len, id)

	return 0
}

func (bfs *BlueFS) mkfs(osdUuid types.UuidD) int {
	log.Debug("osd uuid is %v", osdUuid.UUID)

	initAlloc(bfs)

	initLogger(bfs)

	super := btypes.BlueFsSuperT{
		Version:   uint64(1),
		BlockSize: bfs.bdev[BdevDb].GetBlockSize(),
		OsdUuid:   osdUuid,
		Uuid:      types.GenerateRandomUuid(),
	}

	log.Debug("super uuid is %v", super.Uuid)

	var logFile fileRef
	logFile.fnode.Ino = uint64(1)
	logFile.fnode.PreferBdev = BdevWal

	return 0
}
