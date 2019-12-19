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

type fileRef File
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

type FileReaderBuffer struct {
	blOff       uint64
	pos         uint64
	maxPrefetch uint64
	bl          types.BufferList
}

func CreateFileReaderBuffer(mp uint64) *FileReaderBuffer {
	return &FileReaderBuffer{
		maxPrefetch: mp,
	}
}

func (fb *FileReaderBuffer) getBufEnd() uint64 {
	return fb.blOff + uint64(fb.bl.Length())
}

func (fb *FileReaderBuffer) getBufRemaining(p uint64) uint64 {
	if p > fb.blOff && p < fb.blOff+fb.bl.Length() {
		return fb.blOff + fb.bl.Length() - p
	}
	return 0
}

func (fb *FileReaderBuffer) skip(n uint64) {
	fb.pos += n
}

func (fb *FileReaderBuffer) seek(offset uint64) {
	fb.pos = offset
}

type FileReader struct {
	file      fileRef
	buf       FileReaderBuffer
	random    bool
	ignoreEof bool
}

func CreateFileReader(f fileRef, mpf uint64, rand bool, ie bool) *FileReader {
	fr := &FileReader{
		file: f,
		buf: FileReaderBuffer{
			maxPrefetch: mpf,
		},
		random:    rand,
		ignoreEof: ie,
	}
	fr.file.numReaders++
	return fr
}

type FileLock struct {
	file fileRef
}

func CreateFileLock(_file fileRef) *FileLock {
	return &FileLock{
		file: _file,
	}
}

type BlueFS struct {
	Cct *types.CephContext

	lock    sync.Mutex
	logger  *types.PerfCounters
	dirMap  map[string]dirRef
	fileMap map[uint64]fileRef
	//TODO: dirtyFiles dirty_file_list unknown
	super        types.BlueFsSuperT
	inoLast      uint64
	logSeq       uint64
	logSeqStable uint64
	logWriter    *FileWriter
	logT         types.BlueFsTransactionT
	logFlushing  bool
	logCond      sync.Cond
	newLogJumpTo uint64
	oldLogJumpTo uint64
	newLog       fileRef
	newLogWriter FileWriter

	/* 3 block device
	*	BDEV_DB    db/
	*	BDEV_WAL   db.wal/
	*	BDEV_SLOW  db.slow/
	 */
	bdev           []types.BlockDevice
	ioc            types.IOContext
	blockAll       uint64
	alloc          []*types.Allocator
	AllocSize      []uint64
	pendingRelease []uint64

	slowDevExpander *BlueFSDeviceExpander
}

func CreateBlueFS(cct *types.CephContext) *BlueFS {
	blueFs := &BlueFS{
		Cct: cct,
	}
	return blueFs
}

func (bf *BlueFS) setSlowDeviceExpander(bfe *BlueFSDeviceExpander) {
	bf.slowDevExpander = bfe
}

func (bf *BlueFS) addBlockDevice(deviceId uint8, devPath string) {
	log.Debug("bdev id %d and path %s.", deviceId, devPath)

	types.CreateBlockDevice(bf.Cct, devPath)
}
