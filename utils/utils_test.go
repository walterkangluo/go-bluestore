package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var TestExistDir = "/etc"
var TestNotExistDir = "/etc/hello"

var TestExistFile = "/etc/hosts"
var TestNotExistFile = "/etc/hostname"

func TestExists(t *testing.T) {
	assert := assert.New(t)

	exist, err := PathExists(TestExistDir)
	assert.True(exist)
	assert.Nil(err)

	exist, err = PathExists(TestNotExistDir)
	assert.False(exist)
	assert.Nil(err)

	exist, err = PathExists(TestExistFile)
	assert.True(exist)
	assert.Nil(err)

	exist, err = PathExists(TestNotExistFile)
	assert.False(exist)
	assert.Nil(err)
}

func TestIsDir(t *testing.T) {
	assert := assert.New(t)

	dir := IsDir(TestExistDir)
	assert.True(dir)

	dir = IsDir(TestExistFile)
	assert.False(dir)
}

func TestIsFile(t *testing.T) {
	assert := assert.New(t)

	dir := IsFile(TestExistDir)
	assert.False(dir)

	dir = IsFile(TestExistFile)
	assert.True(dir)
}


