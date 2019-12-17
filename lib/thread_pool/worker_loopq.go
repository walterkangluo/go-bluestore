package thread_pool

import "time"

// implement work by loop-queue which will init queue when create worker array
//  to increase performance, while used more memory resource

type loopQueue struct {
	// worker record
	items  []*worker
	// expiry array
	expiry []*worker
	// queue head
	head   int32
	// queue tail
	tail   int32
	size   int32
	isFull bool
}

func newWorkerLoopQueue(size int32) *loopQueue {
	return &loopQueue{
		items: make([]*worker, size),
		size:  size,
	}
}

func (wq *loopQueue) len() int32 {
	if wq.size == 0 {
		return 0
	}

	if wq.head == wq.tail {
		if wq.isFull {
			return wq.size
		}
		return 0
	}

	if wq.tail > wq.head {
		return wq.tail - wq.head
	}

	return wq.size - wq.head + wq.tail
}

func (wq *loopQueue) isEmpty() bool {
	return wq.head == wq.tail && !wq.isFull
}

func (wq *loopQueue) insert(worker *worker) error {
	if wq.size == 0 {
		return errQueueIsReleased
	}

	if wq.isFull {
		return errQueueIsFull
	}
	wq.items[wq.tail] = worker
	wq.tail++

	if wq.tail == wq.size {
		wq.tail = 0
	}
	if wq.tail == wq.head {
		wq.isFull = true
	}

	return nil
}

func (wq *loopQueue) detach() *worker {
	if wq.isEmpty() {
		return nil
	}

	w := wq.items[wq.head]
	wq.head++
	if wq.head == wq.size {
		wq.head = 0
	}
	wq.isFull = false

	return w
}

func (wq *loopQueue) retrieveExpiry(duration time.Duration) []*worker {
	if wq.isEmpty() {
		return nil
	}

	wq.expiry = wq.expiry[:0]
	expiryTime := time.Now().Add(-duration)

	for !wq.isEmpty() {
		if expiryTime.Before(wq.items[wq.head].recycleTime) {
			break
		}
		wq.expiry = append(wq.expiry, wq.items[wq.head])
		wq.head++
		if wq.head == wq.size {
			wq.head = 0
		}
		wq.isFull = false
	}

	return wq.expiry
}

func (wq *loopQueue) reset() {
	if wq.isEmpty() {
		return
	}

Releasing:
	if w := wq.detach(); w != nil {
		w.task <- nil
		goto Releasing
	}
	wq.items = wq.items[:0]
	wq.size = 0
	wq.head = 0
	wq.tail = 0
}
