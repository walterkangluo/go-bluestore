package types

import "github.com/go-bluestore/common/types"

type BlueStoreIntervalT struct {
	Offset uint64
	Length uint64
}

func CreateBlueStoreIntervalT(o uint64, l uint64) *BlueStoreIntervalT {
	return &BlueStoreIntervalT{
		Offset: o,
		Length: l,
	}
}

type BluesStorePExtentT struct {
	BlueStoreIntervalT
}

func CreateBluesStorePExtentT2(o uint64, l uint64) *BluesStorePExtentT {
	return &BluesStorePExtentT{
		BlueStoreIntervalT{
			Offset: o,
			Length: l,
		},
	}
}

func CreateBluesStorePExtentT1(bs BlueStoreIntervalT) *BluesStorePExtentT {
	return &BluesStorePExtentT{bs}
}

type PExtentVector struct {
	*types.Vector
}
