package types

type BlockDevice struct {
	cct  *CephContext
	path string
}

func CreateBlockDevice(cct *CephContext, path string) *BlockDevice {
	return &BlockDevice{
		cct:  cct,
		path: path,
	}
}

func (*BlockDevice) Open(path string) {

}
