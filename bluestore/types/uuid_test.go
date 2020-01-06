package types

import (
	"fmt"
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

func setNew(aa *string) {
	*aa = "123"
}

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

	setNew(&cc)
	assert.Equal("123", cc)

	//for i := 0; i < len(cc); i++ {
	//	m := cc[i]
	//	m ^= uint8(2)
	//	fmt.Println(m)
	//}
	Netw(&cc, dd, 2)
	fmt.Println(cc)
}

func Netw(data *string, sts string, ll int) {
	//for i := 0; i < ll; i++ {
	//	(*data)[i] ^= sts[i]
	//}
	//*data = "ccc"
	r := []rune(*data)
	r[1] = 'f'
	*data = string(r)
}
