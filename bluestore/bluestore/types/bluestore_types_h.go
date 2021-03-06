package types

import (
	types2 "github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/common/types"
	"time"
)

type BlueStoreIntervalT struct {
	Offset uint64
	Length uint64
}

type BluestoreBdevLabelT struct {
	OsdUUID     types2.UUID
	Size        uint64
	BTime       time.Time
	Description string
	Meta        *types.MapList
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
