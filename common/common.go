package common

import "time"

type ArrayType int

const (
	StackType ArrayType = 1 << iota
	LoopQueueType
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
