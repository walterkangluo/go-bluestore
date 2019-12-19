package bluefs

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
)

func initAlloc(bfs *BlueFS) {
	bfs.AllocSize = make([]uint64, 0, MaxBdev)
}

func (bfs *BlueFS) mkfs(osdUuid types.UuidD) int {
	initAlloc(bfs)
	log.Debug("osd uuid is %v", osdUuid.Data)
	return 0
}
