package thread_pool

import (
	"testing"
	"time"
)

func TestNewLoopQueue(t *testing.T) {
	size := int32(100)
	q := newWorkerLoopQueue(size)
	if q.len() != 0 {
		t.Fatalf("Len error")
	}

	if !q.isEmpty() {
		t.Fatalf("IsEmpty error")
	}

	if q.detach() != nil {
		t.Fatalf("Dequeue error")
	}
}

func TestLoopQueue(t *testing.T) {
	size := int32(10)
	q := newWorkerLoopQueue(size)

	for i := 0; i < 5; i++ {
		err := q.insert(&worker{recycleTime: time.Now()})
		if err != nil {
			break
		}
	}

	if q.len() != 5 {
		t.Fatalf("Len error")
	}

	v := q.detach()
	if nil == v {
		t.Fatalf("Detach error")
	}

	if q.len() != 4 {
		t.Fatalf("Len error")
	}

	time.Sleep(time.Second)

	for i := 0; i < 6; i++ {
		err := q.insert(&worker{recycleTime: time.Now()})
		if err != nil {
			break
		}
	}

	if q.len() != 10 {
		t.Fatalf("Len error")
	}

	err := q.insert(&worker{recycleTime: time.Now()})
	if err == nil {
		t.Fatalf("Enqueue error")
	}

	q.retrieveExpiry(time.Second)

	if q.len() != 6 {
		t.Fatalf("Len error: %d", q.len())
	}
}
