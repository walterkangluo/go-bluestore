package types

type CephContext struct {
	Conf MdConfigT
}

func (cc *CephContext) GetPerfCountersCollection() *PerfCountersCollection {
	return nil
}
