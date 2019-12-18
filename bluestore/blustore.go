package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
)

type SbInfoT struct {
	oidS       []types.GhObject
	sb         types.SharedBlob
	refMap     types.BlueStoreExtentRefMapT
	compressed bool
}
