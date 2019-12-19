package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
)

type SbInfoT struct {
	oidS       []types.GhObject
	sb         SharedBlob
	refMap     BlueStoreExtentRefMapT
	compressed bool
}
