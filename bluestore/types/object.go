package types

type ObjectT struct {
	name string
}

func CreateObject(name string) *ObjectT {
	return &ObjectT{
		name: name,
	}
}

func (ob *ObjectT) Swap(object *ObjectT) {
	ob.name = object.name
}

func (ob *ObjectT) Clean(object *ObjectT) {
	ob.name = *new(string)
}

func (ob *ObjectT) Empty() bool {
	var temp string

	if temp == ob.name {
		return true
	}
	return false
}

type SnapId struct {
	Val uint64
}

func (sp *SnapId) snapIdAdded(val uint64) {
	sp.Val += val
}

func (sp *SnapId) snapIdAutoAdded() {
	sp.snapIdAdded(1)
}
