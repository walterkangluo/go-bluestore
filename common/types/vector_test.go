package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVector_Init(t *testing.T) {
	assert := assert.New(t)

	var vector Vector

	vector.Init()
	assert.Equal(0, vector.Size())
	assert.Equal(0, len(vector.data))
	assert.Equal(false, vector.preAlloc)
}

func TestVector_ReSize(t *testing.T) {
	assert := assert.New(t)
	var v Vector
	v.Init()

	v.PushBack(1)
	v.PushBack(2)
	assert.Equal(2, v.Size())
	assert.Equal(2, len(v.data))
	assert.Equal(false, v.preAlloc)

	v.ReSize(5)
	assert.Equal(2, v.Size())
	assert.Equal(2, len(v.data))

	assert.Equal(2, v.Size())
	assert.Equal(5, v.Capacity())
	assert.Equal(1, v.At(0).(int))
	assert.Equal(2, v.At(1).(int))
}
