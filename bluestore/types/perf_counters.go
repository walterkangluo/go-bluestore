package types

type PerfCounters struct {
	Name        string
	Description string
	Nick        string
	Prio        uint8
}

type PerfCountersCollection struct {
}

func (pc *PerfCountersCollection) Remove(l *PerfCounters) {

}
