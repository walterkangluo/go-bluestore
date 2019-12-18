package hobject

import (
	"github.com/go-bluestore/bluestore/common/object"
	"github.com/go-bluestore/bluestore/common/types"
	"math"
)

type versionT uint32
type genT versionT

var (
	poolMeta      = int64(-1)
	poolTempStart = int64(-2)

	NoGen   = genT(math.MaxUint32)
	NoShard = types.ShardIdT{
		Id: int8(-1),
	}
)

type hObject struct {
	Oid    object.ObjectT
	Snap   object.SnapId
	Pool   int64
	NSpace string

	hash               uint32
	max                bool
	nibbleWiseKeyCache uint32
	hashReverseBits    uint32
	key                string
}

func (ho *hObject) getKey() string {
	return ho.key
}

func (ho *hObject) setKey(_key string) {
	if _key == ho.key {
		ho.key = *new(string)
	} else {
		ho.key = _key
	}
}

func (ho *hObject) getHash() uint32 {
	return ho.hash
}

func isTempPool(pool int64) bool {
	return pool <= poolTempStart
}

func getTempPool(pool int64) int64 {
	return poolTempStart - pool
}

func isMetaPool(pool int64) bool {
	return pool == poolMeta
}

type GhObject struct {
	hObj       hObject
	generation genT
	shareId    types.ShardIdT
	max        bool
}

func CreateGhObjectDefault() *GhObject {
	return &GhObject{
		generation: NoGen,
		shareId:    NoShard,
		max:        false,
	}
}

func CreateGhObject1(object hObject) *GhObject {
	return &GhObject{
		hObj:       object,
		generation: NoGen,
		shareId:    NoShard,
		max:        false,
	}
}

func CreateGhObject3(object hObject, gen genT, shared types.ShardIdT) *GhObject {
	return &GhObject{
		hObj:       object,
		generation: gen,
		shareId:    shared,
		max:        false,
	}
}

func MakePgMeta(pool int64, hash uint32, shared types.ShardIdT) *GhObject {
	h := hObject{
		Oid: *new(object.ObjectT),
		Snap: object.SnapId{
			Val: math.MaxUint64,
		},
		hash: hash,
		max:  false,
		Pool: pool,
	}
	return &GhObject{
		hObj:       h,
		generation: NoGen,
		shareId:    shared,
	}
}

func (gh *GhObject) isPgMeta() bool {
	return gh.hObj.Oid.Empty() && gh.hObj.Pool > 0
}
