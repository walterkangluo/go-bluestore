package types

import (
	"fmt"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

type perfCounterTypeD uint8
type uintT uint8

var (
	PerfCounterNone        = perfCounterTypeD(0)
	PerfCounterTime        = perfCounterTypeD(0x1)
	PerfCounterU64         = perfCounterTypeD(0x2)
	PerfCounterLongRunAvag = perfCounterTypeD(0x4)
	PerfCounterCounter     = perfCounterTypeD(0x8)
	PerfCounterHistogram   = perfCounterTypeD(0x10)

	Bytes = uintT(0)
	None  = uintT(1)

	PrioCritical      = 10
	PrioInteresting   = 8
	PrioUseful        = 5
	PrioUninteresting = 2
	PrioDebugOnly     = 0

	prioDefault = uint8(0)
)

type perfCounterDataAnyD struct {
	name        string
	description string
	nick        string
	prio        uint8
	_type       perfCounterTypeD
	_uint       uintT

	u64       uint64
	avgCount  uint64
	avgCount2 uint64
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
	_uint       uintT
}

func NewPerfCounters(cct *CephContext, name string, lowerBound int, upperBound int) *PerfCounters {
	perf := &PerfCounters{
		mCct:        cct,
		mLowerBound: lowerBound,
		mUpperBound: upperBound,
		Name:        name,
		mLock:       Mutex{},
		mData:       make([]perfCounterDataAnyD, upperBound-lowerBound-1),
	}
	perf.mLock.New("PerfCounters::" + name)
	return perf
}

type perfCounterRef struct {
	Data         *perfCounterDataAnyD
	PerfCounters *PerfCounters
}

type CounterMap map[string]perfCounterRef
type perfCountersSetT []*PerfCounters

func (p perfCountersSetT) add(perf *PerfCounters) {
	var i int
	// sort by name in ascend
	for i = 0; i < len(p); i++ {
		ret := strings.Compare(p[i].Name, perf.Name)
		if ret == 0 {
			log.Warn("%s has exists.", perf.Name)
			return
		}
		if ret == 1 {
			break
		}
	}
	ss := make([]*PerfCounters, len(p)+1)
	rear := append([]*PerfCounters{}, p[i:]...)
	ss = append(p[0:i], perf)
	ss = append(ss, rear...)
	p = ss
}

func (p perfCountersSetT) find(perf *PerfCounters) *PerfCounters {
	for i := 0; i < len(p); i++ {
		ret := strings.Compare(p[i].Name, perf.Name)
		if ret == 0 {
			log.Warn("%s has exists.", perf.Name)
			return p[i]
		}
	}
	return nil
}

func (p perfCountersSetT) erase(perf *PerfCounters) {
	for i := 0; i < len(p); i++ {
		ret := strings.Compare(p[i].Name, perf.Name)
		if ret == 0 {
			log.Warn("%s has exists.", perf.Name)
			p[i] = nil
			return
		}
	}
}

func (p perfCountersSetT) clear(perf *PerfCounters) {
	if 0 == len(p) {
		return
	}
	p = make([]*PerfCounters, 0)
}

type PerfCountersCollection struct {
	mCct     *CephContext
	mLock    sync.Mutex
	mLoggers perfCountersSetT
	byPath   CounterMap
}

func (pc *PerfCountersCollection) Remove(l *PerfCounters) {
	for i := 0; i < len(l.mData); i++ {
		data := l.mData[i]

		path := l.getName() + "." + data.name
		delete(pc.byPath, path)
	}

	perf := pc.mLoggers.find(l)
	utils.AssertTrue(perf != nil)
	pc.mLoggers.erase(l)
}

func (pc *PerfCountersCollection) Clear(l *PerfCounters) {
	pc.mLoggers = make([]*PerfCounters, 0)
	pc.byPath = make(map[string]perfCounterRef)
}

func (pc *PerfCountersCollection) Add(l *PerfCounters) {
	ret := pc.mLoggers.find(l)
	if ret != nil {
		log.Error("exists %s.", l.Name)
		return
	}

	newName := fmt.Sprintf("%s-%d", l.getName(), unsafe.Pointer(l))
	l.setName(newName)

	pc.mLoggers.add(l)

	for i := 0; i < len(l.mData); i++ {
		data := l.mData[i]
		path := l.getName() + "." + data.name
		pc.byPath[path] = perfCounterRef{
			Data:         &data,
			PerfCounters: l,
		}
	}
}

func (pc *PerfCounters) getName() string {
	return pc.Name
}

func (pc *PerfCounters) setName(name string) {
	pc.Name = name
}

func (pc *PerfCounters) Inc(idx int, amt uint64) {
	if !pc.mCct.Conf.Perf {
		return
	}

	utils.AssertTrue(idx > pc.mLowerBound)
	utils.AssertTrue(idx < pc.mUpperBound)

	data := pc.mData[idx-pc.mLowerBound-1]
	if 0 == data._type&PerfCounterU64 {
		return
	}

	if 0 != data._type&PerfCounterLongRunAvag {
		atomic.AddUint64(&data.avgCount, uint64(1))
		atomic.AddUint64(&data.u64, amt)
		atomic.AddUint64(&data.avgCount2, uint64(1))
	} else {
		atomic.AddUint64(&data.u64, amt)
	}
}

type PerCountersBuilder struct {
	mPerfCounters *PerfCounters
	prioDefault   int
}

func CreatePerCountersBuilder(cct *CephContext, name string, lowerBound int, upperBound int) *PerCountersBuilder {
	return &PerCountersBuilder{
		mPerfCounters: NewPerfCounters(cct, name, lowerBound, upperBound),
		prioDefault:   0,
	}
}

func (p *PerCountersBuilder) New(cct *CephContext, name string, lowerBound int, upperBound int) {
	p.mPerfCounters = NewPerfCounters(cct, name, lowerBound, upperBound)
}

func (p *PerCountersBuilder) addImpl(idx int, name string, desc string, nick string, prio int, ty perfCounterTypeD, u uintT) {
	utils.AssertTrue(idx > p.mPerfCounters.mLowerBound)
	utils.AssertTrue(idx < p.mPerfCounters.mUpperBound)

	data := p.mPerfCounters.mData[idx-p.mPerfCounters.mLowerBound-1]
	data.name = name
	data.description = desc
	if len(nick) != 0 {
		utils.AssertTrue(len(nick) < 4)
		data.nick = nick
	}

	data.prio = prioDefault
	if prio != 0 {
		data.prio = uint8(prio)
	}

	data._type = ty
	data._uint = u
	// TODO: add histogram
	// data.histogram
}

func (p *PerCountersBuilder) AddU64Counter(idx int, name string, desc string, nick string, prio int, u int) {
	p.addImpl(idx, name, desc, nick, prio, PerfCounterU64|PerfCounterCounter, uintT(u))
}

func (p *PerCountersBuilder) AddU64(idx int, name string, desc string, nick string, prio int, u int) {
	p.addImpl(idx, name, desc, nick, prio, PerfCounterU64|PerfCounterLongRunAvag, uintT(u))
}

func (p *PerCountersBuilder) AddU64Avg(idx int, name string, desc string, nick string, prio int, u int) {
	p.addImpl(idx, name, desc, nick, prio, PerfCounterU64|PerfCounterLongRunAvag, uintT(u))
}

func (p *PerCountersBuilder) AddTime(idx int, name string, desc string, nick string, prio int, u int) {
	p.addImpl(idx, name, desc, nick, prio, PerfCounterTime, uintT(u))
}

func (p *PerCountersBuilder) AddTimeAvg(idx int, name string, desc string, nick string, prio int, u int) {
	p.addImpl(idx, name, desc, nick, prio, PerfCounterTime|PerfCounterLongRunAvag, uintT(u))
}

func (p *PerCountersBuilder) CreatePerfCounters() *PerfCounters {
	len := len(p.mPerfCounters.mData)
	for i := 0; i < len; i++ {
		data := p.mPerfCounters.mData[i]
		utils.AssertTrue(data._type != PerfCounterNone)
		utils.AssertTrue(0 != data._type&(PerfCounterU64|PerfCounterTime))
	}
	ret := p.mPerfCounters
	p.mPerfCounters = nil
	return ret
}
