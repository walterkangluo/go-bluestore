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

func TestMD5(t *testing.T) {
	assert := assert.New(t)

	var test = "Hello, World!"

	md5 := MD5String(test)
	assert.Equal("65a8e27d8879283831b664bd8b7f0ad4", md5)

	sha1 := SHA1String(test)
	assert.Equal("0a0a9f2a6772942557ab5355d76af442f8f65e01", sha1)

	crc32 := CRC32String(test)
	assert.Equal(uint32(3964322768), crc32)

	testByte := "65a8e27d8879283831b664bd8b7f0ad4"
	toByte := Hex2Bytes(testByte)
	aa := MD5Byte(toByte)
	expectAA := []byte{
		54, 53, 97, 56, 101, 50, 55, 100, 56, 56, 55, 57, 50, 56, 51, 56,
		51, 49, 98, 54, 54, 52, 98, 100, 56, 98, 55, 102, 48, 97, 100, 52,
	}
	assert.Equal(expectAA, aa)

	bb := SHA1Byte(toByte)
	assert.Equal(expectAA, bb)

	cc := CRC32Byte(toByte)
	assert.Equal(uint32(4041457148), cc)
}
