package bluestore

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
)

type SbInfoT struct {
	oidS       []types.GhObject
	sb         SharedBlob
	refMap     BlueStoreExtentRefMapT
	compressed bool
}

func (bs *BlueStore)ReadMeta(key string, value *string) int{
	return 0
}

func (bs *BlueStore)mount(kvOnly bool) int{
	log.Debug("path %s", bs.Path)

	bs.KvOnly = kvOnly

	var mType string
	r := bs.ReadMeta("type", &mType)
	if r < 0 {
		log.Error("expected bluestore, but type is %s", mType)
		return -5
	}

	if mType != "bluestore" {
		log.Error("expected bluestore, but type is %s", mType)
		return -5
	}

	return 0
}

func (bs *BlueStore)Mount() int{
	return bs.mount(false)
}
