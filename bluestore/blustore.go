package bluestore

import (
	"github.com/go-bluestore/bluestore/common/bluestore"
	"github.com/go-bluestore/bluestore/common/hobject"
)

type SbInfoT struct {
	oidS       []hobject.GhObject
	sb         bluestore.SharedBlob
	refMap     bluestore.BlueStoreExtentRefMapT
	compressed bool
}
