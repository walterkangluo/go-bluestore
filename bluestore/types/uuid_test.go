package types

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGenerateRandomUuid(t *testing.T) {
	GenerateRandomUuid()
}

const (
	First = 34300 + iota
	Second
	third
	forth
)

func TestUuidD_IsZero(t *testing.T) {
	assert := assert.New(t)

	var u UUID

	r := u.IsZero()
	assert.True(r)

	u = GenerateRandomUuid()
	r = u.IsZero()
	assert.False(r)

	assert.Equal(First, 34300)
	assert.Equal(Second, 34301)

	aa := new(string)
	assert.NotNil(aa)
	bb := *aa
	assert.Equal("", bb)

	cc := "abc"
	dd := "abcd"

	s := strings.Compare(cc, dd)
	assert.Equal(-1, s)

	s = strings.Compare(dd, cc)
	assert.Equal(1, s)
}
