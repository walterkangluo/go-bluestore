package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
	bs "github.com/go-bluestore/bluestore/types/bluestore"
)

type SbInfoT struct {
	oidS       []types.GhObject
	sb         bs.SharedBlob
	refMap     bs.BlueStoreExtentRefMapT
	compressed bool
}
