package types

import (
	"github.com/go-bluestore/utils"
)

type perfCounterTypeD uint8
type uintT uint8
var (
	PerfCounterNone = perfCounterTypeD(0)
	PerfCounterTime = perfCounterTypeD(0x1)
	PerfCounterU64 = perfCounterTypeD(0x2)
	PerfCounterLongRunAvag = perfCounterTypeD(0x4)
	PerfCounterCounter = perfCounterTypeD(0x8)
	PerfCounterHistogram = perfCounterTypeD(0x10)

	Bytes = uintT(0)
	None = uintT(1)

	PrioCritical = 10
	PrioInteresting = 8
	PrioUseful = 5
	PrioUninteresting = 2
	PrioDebugOnly = 0

	prioDefault = uint8(0)
)

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
	_type       perfCounterTypeD
	_uint    uintT
}

func NewPerfCounters(cct *CephContext, name string, lowerBound int, upperBound int) *PerfCounters{
	perf := &PerfCounters{
		mCct: cct,
		mLowerBound: lowerBound,
		mUpperBound: upperBound,
		Name:name,
		mLock: Mutex{},
		mData: make([]perfCounterDataAnyD, upperBound - lowerBound - 1),
	}
	perf.mLock.New("PerfCounters::" + name)
	return perf
}


type PerfCountersCollection struct {
}

func (pc *PerfCountersCollection) Remove(l *PerfCounters) {

}

func (pc *PerfCounters) Inc(idx int, amt uint64) {
	if !pc.mCct.Conf.Perf {
		return
	}

	utils.AssertTrue(idx > pc.mLowerBound)
	utils.AssertTrue(idx < pc.mUpperBound)

	//data := pc.mData[idx - pc.mLowerBound - 1]
}

type PerCountersBuilder struct {
	mPerfCounters *PerfCounters
	prioDefault int
}

func (p *PerCountersBuilder)New(cct *CephContext, name string, lowerBound int, upperBound int) {
	p.mPerfCounters = NewPerfCounters(cct, name, lowerBound, upperBound)
}

func (p *PerCountersBuilder) addImpl(idx int, name string, desc string, nick string, prio int, ty int, u int) {
	utils.AssertTrue(idx > p.mPerfCounters.mLowerBound)
	utils.AssertTrue(idx < p.mPerfCounters.mUpperBound)

	var data PerfCounters

	data.Name = name
	data.Description = desc
	if len(nick) != 0 {
		utils.AssertTrue(len(nick) < 4)
		data.Nick = nick
	}

	data.Prio = prioDefault
	if prio != 0 {
		data.Prio = uint8(prio)
	}

	data._type = perfCounterTypeD(ty)
	data._uint = uintT(u)

}

func (p *PerCountersBuilder) AddU64Counter(idx int, name string, desc string, nick string, prio int, u int) {

}
