package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElements_GetKey(t *testing.T) {
	assert := assert.New(t)

	fmt.Println("Starting test...")
	ml := NewMapList()
	var a, b, c Keyer
	a = &Elements{"Alice", 1}
	b = &Elements{"Bob", "2"}
	c = &Elements{"Conrad", "3"}
	ml.Push(a)
	ml.Push(b)
	ml.Push(c)
	cb := func(data Keyer) {
		fmt.Println(ml.dataMap[data.GetKey()].Value.(*Elements).key, ml.dataMap[data.GetKey()].Value.(*Elements).val)
	}
	fmt.Println("Print elements in the order of pushing:")
	ml.Walk(cb)

	el, ok := ml.Exists(a)
	assert.Nil(ok)
	assert.Equal(a.GetKey(), el.(*Elements).GetKey())
	assert.Equal(a.GetVal(), el.(*Elements).GetVal())
	fmt.Printf("Size of MapList: %d \n", ml.Size())
	ml.Remove(b)
	fmt.Println("After removing b:")
	ml.Walk(cb)
	fmt.Printf("Size of MapList: %d \n", ml.Size())
}
