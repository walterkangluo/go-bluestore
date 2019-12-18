package bluestore

import (
	"github.com/go-bluestore/bluestore/common"
	"github.com/go-bluestore/log"
)

var (
	BDEV_WAL  = 0
	BDEV_DB   = 1
	BDEV_SLOW = 2
	MAX_BDEV  = 3
)

type BlueFS struct {
	cct             *common.CephContext
	slowDevExpander common.BlueFSDeviceExpander
	bdev            map[uint8]BlockDevice
}

func CreateBlueFS(cct *common.CephContext) *BlueFS {
	blueFs := &BlueFS{
		cct: cct,
	}
	return blueFs
}

func (bf *BlueFS) setSlowDeviceExpander(bfe common.BlueFSDeviceExpander) {
	bf.slowDevExpander = bfe
}

func (bf *BlueFS) addBlockDevice(deviceId uint8, devPath string) {
	log.Debug("bdev id %d and path %s.", deviceId, devPath)

	CreateBlockDevice(bf.cct, devPath)
}
