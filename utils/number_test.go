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

	a := ROUND_UP_TO(5, 2)
	assert.Equal(int64(6), a)

	a = ROUND_UP_TO(5, 4)
	assert.Equal(int64(8), a)
}


func TestDIV_ROUND_UP(t *testing.T) {
	assert := assert.New(t)

	a := DIV_ROUND_UP(5, 2)
	assert.Equal(int64(3), a)
}