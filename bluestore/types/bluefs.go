package types

import (
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
	cct             *CephContext
	slowDevExpander BlueFSDeviceExpander
	bdev            map[uint8]BlockDevice
}

func CreateBlueFS(cct *CephContext) *BlueFS {
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

	CreateBlockDevice(bf.cct, devPath)
}
