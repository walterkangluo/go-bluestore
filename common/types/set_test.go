package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_Init(t *testing.T) {
	assert := assert.New(t)

	var s Set

	s.Init()
	assert.Equal(uint32(0), s.Size())
	assert.Nil(s.Begin())
	assert.Nil(s.Back())
}

func TestSet_Push(t *testing.T) {
	assert := assert.New(t)

	var s Set

	s.Init()
	assert.Equal(uint32(0), s.Size())
	assert.Nil(s.Begin())
	assert.Nil(s.Back())

	var ss S

	ss = make(map[interface{}]interface{})
	ss["1"] = "first"
	s.Push(ss)
	assert.Equal(uint32(1), s.Size())

	cc := s.Begin()
	val, ok := cc["1"]
	assert.True(ok)
	assert.Equal("first", val)
}
