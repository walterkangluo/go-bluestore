package bluefs

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
	"sync"
)

const (
	BdevWal  = 0
	BdevDb   = 1
	BdevSlow = 2
	MaxBdev  = 3

	WriteUnknown = 0
	WriteWal     = 1
	WriteSst     = 2
)

type BlueFSDeviceExpander struct {
}

type File struct {
	types.RefCountedObject

	fnode    types.BlueFsFnodeT
	refs     int
	dirtySeq uint64
	locked   bool
	deleted  bool
	// dirtyItem []

	numReaders int
	numWriters int
	numReading int
}

type fileRef []File
type Dir struct {
	types.RefCountedObject
	fileMap map[string]fileRef
}

type dirRef []Dir
type FileWriter struct {
	file fileRef
	pos  uint64
	// buff bufferList
	// tailblock bufferList
	// bufferAppender bufferlist
	writerType int
	lock       sync.Mutex
	//iocv [MaxBdev]types.AioContext
}

type BlueFS struct {
	cct *types.CephContext

	slowDevExpander BlueFSDeviceExpander
	bdev            map[uint8]types.BlockDevice
}

func CreateBlueFS(cct *types.CephContext) *BlueFS {
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

	types.CreateBlockDevice(bf.cct, devPath)
}
