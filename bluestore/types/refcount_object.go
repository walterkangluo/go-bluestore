package types

type RefCountedObject struct {
	NRef uint64
	Cct  *CephContext
}
