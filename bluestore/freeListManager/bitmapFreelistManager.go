package freeListManager

import (
	"github.com/go-bluestore/bluestore/kv"
	"github.com/go-bluestore/bluestore/kv/common"
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"sync"
)

type BitmapFreelistManager struct {
	*FreelistManager
	metaPrefix   string
	bitmapPrefix string
	kvdb         *kv.KeyValueDB
	Lock         sync.Mutex

	size          uint64
	bytesPerBlock uint64
	blocksPerKey  uint64
	bytesPerKey   uint64
	blocks        uint64
	keyMask       uint64

	allSetBl        types.BufferList
	enumerateOffset uint64
	enumerateBlPos  int
}

func NewBitmapFreelistManager(cct *types.CephContext, db *kv.KeyValueDB,
	metaPrefix string, bitmapPrefix string) (bmfm *BitmapFreelistManager) {

	bmfm = &BitmapFreelistManager{
		metaPrefix:     metaPrefix,
		kvdb:           db,
		enumerateBlPos: 0,
		bitmapPrefix:   bitmapPrefix,
	}
	bmfm.FreelistManager.New(cct)

	return
}

func (bf *BitmapFreelistManager) Create(newSize uint64, graunlarity uint64, txn common.Transaction) {
	bf.bytesPerBlock = graunlarity
	utils.AssertTrue(utils.ISP2(bf.bytesPerBlock))

	size := utils.P2Align(newSize, bf.bytesPerBlock)
	bf.blocksPerKey = bf.Cct.Conf.BlueStoreFreelistBlocksPerKey

	// initMisc

	bf.blocks = size / bf.bytesPerBlock
	if bf.blocks != bf.blocksPerKey*bf.blocksPerKey {
		bf.blocks = (bf.blocks/bf.blocksPerKey + 1) * bf.blocksPerKey
		log.Debug("rouding blocks up from %x to %x, (%x blocks)", size, bf.blocks*bf.bytesPerBlock, bf.blocks)
	}
	// TODO: Add more
}

type XorMergeOperator struct {
}

func (xo *XorMergeOperator) MergeNonexistent(rData string, rLen int, newValue *string) {
	*newValue = rData[:rLen]
}

func (xo *XorMergeOperator) Merge(rData string, rLen int, lData string, lLen int, newValue *string) {
	utils.AssertTrue(rLen == lLen)

	r := []rune(*newValue)
	s := []rune(rData)
	for i := 0; i < rLen; i++ {
		r[i] ^= s[i]
	}
	*newValue = string(r)
}

func (xo *XorMergeOperator) Name() string {
	return "bitwise_xor"
}

func (fm *BitmapFreelistManager) SetupMergeOperators(db *kv.KeyValueDB, prefix string) {
	merOp := new(XorMergeOperator)
	fm.kvdb.SetMergeOperator(prefix, merOp)
}
