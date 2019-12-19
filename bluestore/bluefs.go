package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/bluestore/types/bluefs"
	"github.com/go-bluestore/log"
)

func initAlloc(bfs *bluefs.BlueFS) {
	bfs.AllocSize = make([]uint64, 0, bluefs.MaxBdev)
}

func mkfs(osdUuid types.UuidD) int {
	bfs := bluefs.CreateBlueFS(nil)
	initAlloc(bfs)
	log.Debug("osd uuid is %v", osdUuid.Data)
	return 0
}
