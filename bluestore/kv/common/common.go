package common

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/utils"
)

//type Transaction TransactionImpl

type Transaction struct {
	TransactionImpl
}

type TransactionImpl struct {
}

func (ti *TransactionImpl) Set1(prefix string, toSet map[string]types.BufferList) {
	for key, value := range toSet {
		ti.Set3(prefix, key, value)
	}
}

func (ti *TransactionImpl) Set2(prefix string, toSetBl types.BufferList) {
}

func (ti *TransactionImpl) Set3(prefix string, k string, bl types.BufferList) {

}

func (ti *TransactionImpl) Set4(prefix string, toSetBl types.BufferList) {
	p := toSetBl.Begin()
	var num uint32
	toSetBl.Decode(utils.Int32ToBytes(num), p)
	for {
		var key string
		//var value types.BufferList
		toSetBl.Decode([]byte(key), p)
		//toSetBl.Decode(value, p)
	}
}
