package aio

import (
	"unsafe"
)

type iocb struct {
}

type AIOT struct {
	iocb   iocb
	priv   unsafe.Pointer
	fd     int
	length uint64
	offset uint64
	rval   int64
}

func CreateAIOT(p unsafe.Pointer, f int) *AIOT {
	return &AIOT{
		priv: p,
		fd:   f,
		rval: int64(-1000),
	}
}

func (aio *AIOT) PWriteV(_offset uint64, _length uint64) {
	aio.length = _length
	aio.offset = _offset
	//TODO: io_prep_pritev
}

func (aio *AIOT) PRead(_offset uint64, _length uint64) {
	aio.length = _length
	aio.offset = _offset
	// TODO: io_prep_read
}

func (aio *AIOT) GetReturnValue() int64 {
	return aio.rval
}

// TODO: replace by io_context_t
type ioContextT int
type AioQueueT struct {
	maxIODepth int
	ctx        ioContextT
	aioIter    AIOT
}

func Create(maxIoDepth int) AioQueueT {
	return AioQueueT{
		maxIODepth: maxIoDepth,
		ctx:        0,
	}
}
