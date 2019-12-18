package bluefs

import (
	bdv "github.com/go-bluestore/bluestore/common/blockdevice"
	ctx "github.com/go-bluestore/bluestore/common/ceph_context"
	"github.com/go-bluestore/log"
)

var (
	BdevWal  = 0
	BdevDb   = 1
	BdevSlow = 2
	MaxBdev  = 3
)

type BlueFSDeviceExpander struct {
}

type BlueFS struct {
	cct             *ctx.CephContext
	slowDevExpander BlueFSDeviceExpander
	bdev            map[uint8]bdv.BlockDevice
}

func CreateBlueFS(cct *ctx.CephContext) *BlueFS {
	blueFs := &BlueFS{
		cct: cct,
	}
	return blueFs
}

func (bf *BlueFS) setSlowDeviceExpander(bfe BlueFSDeviceExpander) {
	bf.slowDevExpander = bfe
}

func (bf *BlueFS) addBlockDevice(deviceId uint8, devPath string) {
	log.Debug("bdev id %d and path %s.", deviceId, devPath)

	bdv.CreateBlockDevice(bf.cct, devPath)
}
