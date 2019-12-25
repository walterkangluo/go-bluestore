package types

import (
	"github.com/go-bluestore/common"
)

type ObjectStore struct {
	Cct *CephContext

	Path string
}

func (obs *ObjectStore) ReadMeta(key string, value *string) int {
	var buf string
	r := common.SafeReadFile(obs.Path, key, &buf, 4096)
	if r <= 0 {
		return int(r)
	}
	var i int
	for i = len(buf) - 1; i >= 0; i-- {
		if buf[i] != ' ' {
			break
		}
	}
	*value = buf[:i]

	return 0
}

func (obs *ObjectStore) WriteMeta(key string, value string) int {
	return 0
}
