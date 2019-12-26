package types

import "github.com/go-bluestore/utils"

type perfCounterDataAnyD struct {
	name        string
	description string
	nick        string
	prio        uint
}

type perfCounterDataVecT []perfCounterDataAnyD

type PerfCounters struct {
	Name        string
	Description string
	Nick        string
	Prio        uint8

	// private
	mCct        *CephContext
	mLowerBound int
	mUpperBound int
	mName       string
	mLockName   string
	prioAdjust  int
	mLock       Mutex
	mData       perfCounterDataVecT
}

type PerfCountersCollection struct {
}

func (pc *PerfCountersCollection) Remove(l *PerfCounters) {

}

var (
	PertCounterNone       = 0
	PertCounterTime       = 0x1
	PertCounterU64        = 0x2
	PertCounterLongRunAvg = 0x4
	PertCounterCounter    = 0x8
	PertCounterHistogram  = 0x10
)

func (pc *PerfCounters) Inc(idx int, amt uint64) {
	if !pc.mCct.Conf.Perf {
		return
	}

	utils.AssertTrue(idx > pc.mLowerBound)
	utils.AssertTrue(idx < pc.mUpperBound)

	//data := pc.mData[idx - pc.mLowerBound - 1]

}
