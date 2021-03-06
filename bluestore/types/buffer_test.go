package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBufferList_Add(t *testing.T) {
	var bl BufferList
	bl.Init()
	assert.Equal(t, uint64(0), bl.size)
	data := []byte("abcdef")

	bl.Add(data, uint64(len(data)))

	assert.Equal(t, uint64(1), bl.size)
}
