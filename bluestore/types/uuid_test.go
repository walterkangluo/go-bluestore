package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRandomUuid(t *testing.T) {
	GenerateRandomUuid()
}

func TestUuidD_IsZero(t *testing.T) {
	assert := assert.New(t)

	var u UuidD

	r := u.IsZero()
	assert.True(r)

	u = GenerateRandomUuid()
	r = u.IsZero()
	assert.False(r)
}
