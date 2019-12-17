package thread_pool

import (
	"errors"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/lib"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (

	// ErrInvalidPoolSize will be returned when setting a negative number as pool capacity.
	ErrInvalidPoolSize = errors.New("invalid size for pool")

	// ErrInvalidPoolExpiry will be returned when setting a negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrExistsPoolName will be return when create pool by exist pool name
	ErrPoolNameExist = errors.New("pool has exist with given name")

	// CLOSED represents that the pool is closed.
	CLOSED = int32(1)

	// workerChanCap determines whether the channel of a worker should be a buffered channel
	// to get the best performance. Inspired by fasthttp at https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	WorkerChanCap = func() int {
		// Use blocking workerChan if GOMAXPROCS=1.
		// This immediately switches Serve to WorkerFunc, which results
		// in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the Serve caller (Acceptor) may lag accepting
		// new connections if WorkerFunc is CPU-bound.
		return 1
	}()

	PoolRecords = make(map[string]*Pool, 0)
)

// Pool accepts the tasks from client
// it limits the total of goroutines to a given number by recycling goroutines.
type Pool struct {
	// name of the pool
	name string

	// capacity of the pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// workers is a slice that store the available workers.
	workers workerArray

	// release is used to notice the pool to closed itself.
	release int32

	// lock for synchronous operation.
	lock sync.Locker

	// cond for waiting to get a idle worker.
	cond *sync.Cond

	// once makes sure releasing this pool will just be done for one time.
	once sync.Once

	// workerCache speeds up the obtainment of the an usable worker in function:retrieveWorker.
	workerCache sync.Pool

	// blockingNum is the number of the goroutines already been blocked on pool.Submit, protected by pool.lock
	blockingNum int

	//	options *Options
	flags common.PoolFlags
}

func NewThreadPool(name string, size int32, flags common.PoolFlags) (*Pool, error) {

	if size <= 0 {
		return nil, ErrInvalidPoolSize
	}

	enableExpiry := true
	if expiry := flags.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		enableExpiry = false
	}

	if _, exist := PoolRecords[name]; exist {
		return nil, ErrPoolNameExist
	}

	pool := &Pool{
		name:     name,
		capacity: size,
		lock:     lib.NewSpinLock(),
		flags:    flags,
	}

	pool.workerCache.New = func() interface{} {
		return &worker{
			pool: pool,
			task: make(chan func(), WorkerChanCap),
		}
	}

	if flags.PreAlloc {
		pool.workers = newWorkerArray(common.LoopQueueType, size)
	} else {
		pool.workers = newWorkerArray(common.StackType, 0)
	}

	pool.cond = sync.NewCond(pool.lock)

	if enableExpiry {
		// Start a goroutine to clean up expired workers periodically.
		go pool.periodicallyPurge()
	}

	// Add to record
	PoolRecords[name] = pool
	return pool, nil
}

// Clear expired workers periodically.
func (p *Pool) periodicallyPurge() {
	heartbeat := time.NewTicker(p.flags.ExpiryDuration)
	defer heartbeat.Stop()

	// will issue interval of p.flags.ExpiryDuration
	for range heartbeat.C {
		if atomic.LoadInt32(&p.release) == CLOSED {
			break
		}

		p.lock.Lock()
		expiredWorkers := p.workers.retrieveExpiry(p.flags.ExpiryDuration)
		p.lock.Unlock()

		// Notify obsolete workers to stop.
		// This notification must be outside the p.lock, since w.task
		// may be blocking and may consume a lot of time if many workers
		// are located on non-local CPUs.
		for i := range expiredWorkers {
			expiredWorkers[i].task <- nil
		}

		// There might be a situation that all workers have been cleaned up(no any worker is running)
		// while some invokers still get stuck in "p.cond.Wait()",
		// then it ought to wakes all those invokers.
		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}

// Running returns the number of the currently running goroutines.
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// Free returns the available goroutines to work.
func (p *Pool) Free() int {
	return p.Cap() - p.Running()
}

// Cap returns the capacity of this pool.
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// Tune changes the capacity of this pool.
func (p *Pool) Tune(size int) {
	if size < 0 || p.Cap() == size || p.flags.PreAlloc {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
}

// Release Closes this pool.
func (p *Pool) Release() {
	p.once.Do(func() {
		atomic.StoreInt32(&p.release, 1)
		p.lock.Lock()
		p.workers.reset()
		p.lock.Unlock()
		delete(PoolRecords, p.name)
	})
}

// Submit submits a task to this pool.
func (p *Pool) Submit(task func()) error {
	if atomic.LoadInt32(&p.release) == CLOSED {
		return ErrPoolClosed
	}
	var w *worker
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.task <- task
	return nil
}

// retrieveWorker returns a available worker to run the tasks.
func (p *Pool) retrieveWorker() *worker {
	var w *worker

	// to get available worker and run the tasks
	spawnWorker := func() {
		w = p.workerCache.Get().(*worker)
		w.run()
	}

	p.lock.Lock()

	// try to get a available work
	w = p.workers.detach()
	if w != nil {
		// success
		p.lock.Unlock()
	} else if p.Running() < p.Cap() {
		// no available worker, spawn a worker
		p.lock.Unlock()
		// spawn and run
		spawnWorker()
	} else {
		if p.flags.NonBlocking {
			// pool not allow appear blocking, so return nil
			p.lock.Unlock()
			return nil
		}

	Reentry:
		// pool allow blocking with MaxBlockingTasks
		if p.flags.MaxBlockingTasks != 0 && p.blockingNum >= p.flags.MaxBlockingTasks {
			p.lock.Unlock()
			return nil
		}

		// record blocking and wait be schedule
		p.blockingNum++
		p.cond.Wait()
		p.blockingNum--
		if p.Running() == 0 {
			p.lock.Unlock()
			spawnWorker()
			return w
		}

		w = p.workers.detach()
		if w == nil {
			goto Reentry
		}

		p.lock.Unlock()
	}
	return w
}

// incRunning increases the number of the currently running goroutines.
func (p *Pool) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

// decRunning decreases the number of the currently running goroutines.
func (p *Pool) decRunning() {
	atomic.AddInt32(&p.running, -1)
}

// revertWorker puts a worker back into free pool, recycling the goroutines.
func (p *Pool) revertWorker(worker *worker) bool {
	if atomic.LoadInt32(&p.release) == CLOSED || p.Running() > p.Cap() {
		return false
	}

	worker.recycleTime = time.Now()
	p.lock.Lock()

	err := p.workers.insert(worker)
	if err != nil {
		return false
	}

	// Notify the invoker stuck in 'retrieveWorker()' of there is an available worker in the worker queue.
	p.cond.Signal()
	p.lock.Unlock()
	return true
}


