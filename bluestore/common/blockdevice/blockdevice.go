package blockdevice

import (
	ctx "github.com/go-bluestore/bluestore/common/ceph_context"
)

type BlockDevice struct {
	cct  *ctx.CephContext
	path string
}

func CreateBlockDevice(cct *ctx.CephContext, path string) *BlockDevice {
	return &BlockDevice{
		cct:  cct,
		path: path,
	}
}

func (*BlockDevice) Open(path string) {

}
