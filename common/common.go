package common

import (
	"io"
	"os"
	"syscall"
	"time"
)

type ArrayType int

const (
	StackType ArrayType = 1 << iota
	LoopQueueType
)

var (
	SPDKPrefix = "spdk:"
)

// Options contains all options which will be applied when instantiating a ants pool.
type PoolFlags struct {
	// ExpiryDuration set the expired time of every worker.
	ExpiryDuration time.Duration

	// PreAlloc indicate whether to make memory pre-allocation when initializing Pool.
	PreAlloc bool

	// Max number of goroutine blocking on pool.Submit.
	// 0 (default value) means no such limit.
	MaxBlockingTasks int

	// When NonBlocking is true, Pool.Submit will never be blocked.
	// ErrPoolOverload will be returned when Pool.Submit cannot be done at once.
	// When NonBlocking is true, MaxBlockingTasks is inoperative.
	NonBlocking bool

	// PanicHandler is used to handle panics from each worker goroutine.
	// if nil, panics will be thrown out again from worker goroutines.
	PanicHandler func(interface{})
}

func SafeRead(f *os.File, buf *string, count int64) (int64, error) {
	var cnt int64
	var data = make([]byte, count)
	var err error
	for cnt < count {
		r, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				return 0, err
			}
			if err == syscall.EINTR {
				continue
			}
			return 0, err
		}
		cnt += int64(r)
		*buf += string(data)
	}
	return cnt, err
}

func SafeReadFile(base string, file string, val *string, valLen int64) int64 {
	fn := base + file
	f, openErr := os.OpenFile(fn, syscall.O_RDONLY, 0)
	if openErr != nil {
		return -1
	}
	len, readErr := SafeRead(f, val, valLen)
	if readErr != nil {
		for closeErr := f.Close(); closeErr != nil && syscall.EINTR == closeErr; closeErr = f.Close() {
		}
		return -2
	}

	for closeErr := f.Close(); closeErr != nil && syscall.EINTR == closeErr; closeErr = f.Close() {
	}

	return len
}
