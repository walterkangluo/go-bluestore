package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBufferList_Length(t *testing.T) {
	var bf BufferList

	bf.Init()
	assert.Equal(t, bf.Length(), uint64(0))
}
