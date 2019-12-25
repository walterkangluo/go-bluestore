package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestISP2(t *testing.T) {
	assert := assert.New(t)

	assert.True(ISP2(4))
	assert.False(ISP2(6))
}

func TestROUND_UP_TO(t *testing.T) {
	assert := assert.New(t)

	a := RoundUpTo(5, 2)
	assert.Equal(int64(6), a)

	a = RoundUpTo(5, 4)
	assert.Equal(int64(8), a)
}

func TestDIV_ROUND_UP(t *testing.T) {
	assert := assert.New(t)

	a := DivRoundUp(5, 2)
	assert.Equal(int64(3), a)
}

func TestCtx(t *testing.T) {
	assert := assert.New(t)

	a := Ctx(4)
	assert.Equal(a, uint(2))

	a = Ctx(5)
	assert.Equal(a, uint(3))

	a = Ctx(6)
	assert.Equal(a, uint(3))

	a = Ctx(8)
	assert.Equal(a, uint(3))

	a = Ctx(9)
	assert.Equal(a, uint(4))

}
