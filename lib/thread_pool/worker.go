package thread_pool

import (
	"errors"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/log"
	"runtime"
	"time"
)

var (
	// errQueueIsFull will be returned when the worker queue is full.
	errQueueIsFull = errors.New("the queue is full")

	// errQueueIsReleased will be returned when trying to insert item to a released worker queue.
	errQueueIsReleased = errors.New("the queue length is zero")

	// ErrPoolClosed will be returned when submitting task to a closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")
)

type workerArray interface {
	len() int32
	isEmpty() bool
	insert(worker *worker) error
	detach() *worker
	retrieveExpiry(duration time.Duration) []*worker
	reset()
}

// worker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type worker struct {
	// pool who owns this worker.
	pool *Pool

	// task is a job should be done.
	task chan func()

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

func newWorkerArray(aType common.ArrayType, size int32) workerArray {
	switch aType {
	case common.StackType:
		return newWorkerStack(size)
	case common.LoopQueueType:
		return newWorkerLoopQueue(size)
	default:
		return newWorkerStack(size)
	}
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *worker) run() {
	w.pool.incRunning()
	go func() {
		defer func() {
			w.pool.decRunning()
			w.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				if ph := w.pool.flags.PanicHandler; ph != nil {
					ph(p)
				} else {
					log.Warn("worker exits from a panic: %v.", p)
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					log.Warn("worker exits from panic: %s.", string(buf[:n]))
				}
			}
		}()

		for f := range w.task {
			if f == nil {
				return
			}
			f()
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}
