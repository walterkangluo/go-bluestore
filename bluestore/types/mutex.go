package types

import (
	"sync"
)

type Mutex struct {
	// private
	name      string
	id        int
	recursive bool
	lockDep   bool
	backtrace bool

	m      sync.Mutex
	nlock  int
	lockBy sync.Mutex
	cct    *CephContext
	logger *PerfCounters
}

func (mu *Mutex) New(name string, a ...interface{}) {
	mu.name = name

	if len(a) == 0 {
		mu.recursive = false
		mu.lockDep = true
		mu.backtrace = false
		mu.cct = nil
	}
}
