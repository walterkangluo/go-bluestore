package bluefs

import (
	"github.com/go-bluestore/bluestore/blockdevice"
	btypes "github.com/go-bluestore/bluestore/bluefs/types"
	"github.com/go-bluestore/bluestore/types"
	ctypes "github.com/go-bluestore/common/types"
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

	lBlueFsFirst     = 732600
	lBlueFsGiftBytes = 732601
)

type BlueFSDeviceExpander struct {
}

type File struct {
	types.RefCountedObject

	fnode    btypes.BlueFsFnodeT
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

func (d *Dir) new() {
	d.RefCountedObject.New(nil, 0)
}

func (d *Dir) empty() bool {
	if len(d.fileMap) == 0 {
		return true
	}
	return false
}

type dirRef struct {
	*Dir
}

type FileWriter struct {
	file fileRef
	pos  uint64
	// buff bufferList
	// tailblock bufferList
	// bufferAppender bufferlist
	writerType int
	lock       sync.Mutex
	iocv       [MaxBdev]*blockdevice.IOContext
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

type blockInfo struct {
	start uint64
	len   uint64
}

func (bi blockInfo) GetStart() uint64 {
	return bi.start
}

func (bi blockInfo) GetLen() uint64 {
	return bi.len
}

//type blockInfoList []blockInfo

type blockInfoList struct {
	segment []blockInfo
	size    uint32
}

func (bl *blockInfoList) insert(offset uint64, length uint64) {
	b := blockInfo{
		start: offset,
		len:   length,
	}
	bl.segment = append(bl.segment, b)
	bl.size++
}

func (bl *blockInfoList) empty() bool {
	return bl.size == 0
}

type BlueFS struct {
	Cct *types.CephContext

	lock    sync.Mutex
	logger  *types.PerfCounters
	dirMap  map[string]dirRef
	fileMap map[uint64]fileRef

	//TODO: dirtyFiles dirty_file_list unknown
	dirtyFiles   map[uint64]interface{}
	super        *btypes.BlueFsSuperT
	inoLast      uint64
	logSeq       uint64
	logSeqStable uint64
	logWriter    *FileWriter
	logT         btypes.BlueFsTransactionT
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
	bdev           *ctypes.Vector  // *types.BlockDevice
	ioc            *ctypes.Vector  // types.IOContext
	blockAll       []blockInfoList // []blockInfoList
	blockTotal     []uint64        // []uint64
	alloc          *ctypes.Vector  // []allocator.Allocator
	allocSize      *ctypes.Vector  // []uint64
	pendingRelease *ctypes.Vector  // []uint64

	slowDevExpander *BlueFSDeviceExpander
}

func CreateBlueFS(cct *types.CephContext) (blueFs *BlueFS) {
	blueFs.Cct = cct

	blueFs.blockAll = make([]blockInfoList, 0)
	blueFs.bdev.Init()
	blueFs.ioc.Init()
	blueFs.blockTotal = make([]uint64, 0)
	blueFs.alloc.Init()
	blueFs.allocSize.Init()
	blueFs.pendingRelease.Init()
	return blueFs
}

func getSuperOffset() uint64 {
	return 4096
}

func getSuperLength() uint64 {
	return 4096
}
