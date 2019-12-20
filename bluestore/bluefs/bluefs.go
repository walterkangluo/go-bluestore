package bluefs

import (
	al "github.com/go-bluestore/bluestore/allocator"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
)

func (bfs *BlueFS) initAlloc() {
	bfs.alloc = make([]*al.Allocator, MaxBdev)
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

	var buleFsFile = []string{"bluefs-wal", "bluefs-db", "bluefs-slow"}
	for id := 0; id < len(bfs.bdev); id++ {
		if nil == bfs.bdev[id] {
			continue
		}

		utils.AssertTrue(bfs.bdev[id].GetSize() > 0)
		utils.AssertTrue(bfs.allocSize[id] > 0)

		log.Debug("bdev id %d, allocSize %d, size %d.", buleFsFile[id], bfs.allocSize[id], bfs.bdev[id].GetSize())

		//bfs.alloc[id] =

	}
}

func (bfs *BlueFS) mkfs(osdUuid types.UuidD) int {
	bfs.initAlloc()
	log.Debug("osd uuid is %v", osdUuid.Data)
	return 0
}
