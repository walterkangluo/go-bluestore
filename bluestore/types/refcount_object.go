package types

type RefCountedObject struct {
	NRef uint64
	Cct  *CephContext
}

func (rc *RefCountedObject) New(c *CephContext, n uint64) {
	rc.Cct = c
	rc.NRef = n
}

func CreateRefCountedObject(c *CephContext, n uint64) *RefCountedObject {
	return &RefCountedObject{
		Cct:  c,
		NRef: n,
	}
}
