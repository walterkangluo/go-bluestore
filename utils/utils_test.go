package utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestHex2Bytes(t *testing.T) {
	assert := assert.New(t)
	type Server struct {
		ServerName string `json:"servername"`
		ServerIP   int    `json:"serverip"`
	}

	type Serverslice struct {
		Servers []Server `json:"servers"`
	}

	var s Serverslice
	//func append(slice []Type, elems ...Type) []Type
	s.Servers = append(s.Servers, Server{ServerName: "Beijing", ServerIP: 100})
	s.Servers = append(s.Servers, Server{ServerName: "Xi'an", ServerIP: 200})
	//slice里面嵌套结构体[{},{}] 遍历出来的是slice里面包含json串

	b, err := json.Marshal(s)
	assert.Nil(err)

	var ss Serverslice
	//将json的字符串转换成s对象(这里用的是指针的方式，所以可以直接修改底层结构体中的数据)
	//需要读取的json字符串必须先写入[]byte类型的对象中(二进制对象文件)
	err = json.Unmarshal(b, &ss)
	assert.Nil(err)
	assert.Equal(2, len(ss.Servers))
	assert.Equal(s.Servers[0].ServerIP, ss.Servers[0].ServerIP)
	assert.Equal(s.Servers[1].ServerIP, ss.Servers[1].ServerIP)

	assert.Equal(s.Servers[0].ServerName, ss.Servers[0].ServerName)
	assert.Equal(s.Servers[1].ServerName, ss.Servers[1].ServerName)
}

func TestEnsureFolderExist(t *testing.T) {
	var aa []int
	type H struct {
		a int
	}
	var mm []H
	s := len(aa)
	n := len(mm)
}

