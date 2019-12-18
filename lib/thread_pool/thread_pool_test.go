package thread_pool

import (
	"github.com/go-bluestore/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


var DefaultPoolName = "test_pool"

func TestNewThreadPool(t *testing.T) {
	assert := assert.New(t)

	flags := common.PoolFlags{}
	size := int32(-1)
	pool, err := NewThreadPool(DefaultPoolName, size, flags)
	assert.Error(ErrInvalidPoolSize, err)
	assert.Nil(pool)

	size = int32(10)
	flags.ExpiryDuration = time.Duration(-1)
	pool, err = NewThreadPool(DefaultPoolName, size, flags)
	assert.Error(ErrInvalidPoolExpiry, err)
	assert.Nil(pool)

	flags.ExpiryDuration = time.Duration(0)
	pool, err = NewThreadPool(DefaultPoolName, size, flags)
	assert.Nil(err)
	assert.NotNil(pool)

	_, err = NewThreadPool(DefaultPoolName, size, flags)
	assert.Error(ErrPoolNameExist, err)

	assert.Equal(int32(0), pool.release)
	pool.Release()
	assert.Equal(int32(1), pool.release)

	pool, err = NewThreadPool(DefaultPoolName, size, flags)
	assert.Nil(err)
	assert.NotNil(pool)
	pool.Release()
}

func TestPool_Cap(t *testing.T) {
	assert := assert.New(t)

	size := int32(10)
	flags := common.PoolFlags{}
	pool, err := NewThreadPool(DefaultPoolName, size, flags)
	assert.Nil(err)
	assert.NotNil(pool)

	assert.Equal(int(size), pool.Cap())
	assert.Equal(0, pool.Running())
	assert.Equal(int(size), pool.Free())

	pool.Tune(-1)
	assert.Equal(int(size), pool.Cap())

	pool.Tune(5)
	assert.Equal(5, pool.Cap())
	pool.Release()
}

func TestPool_Submit(t *testing.T) {
	assert := assert.New(t)

	size := int32(10)
	flags := common.PoolFlags{}
	pool, err := NewThreadPool(DefaultPoolName, size, flags)
	assert.Nil(err)
	assert.NotNil(pool)

	var com = make(chan int)
	err = pool.Submit(func() {
		com <- 2
		return
	})
	ch := <- com
	assert.Equal(2, ch)

	pool.Release()
	err = pool.Submit(func() {
		com <- 3
		return
	})
	assert.Error(ErrPoolClosed, err)
}