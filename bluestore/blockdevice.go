package bluestore

import "github.com/go-bluestore/bluestore/common"

type BlockDevice struct {
	cct  *common.CephContext
	path string
}

func CreateBlockDevice(cct *common.CephContext, path string) *BlockDevice {
	return &BlockDevice{
		cct:  cct,
		path: path,
	}
}

func (*BlockDevice) Open(path string) {

}
