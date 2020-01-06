package env

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import (
	"github.com/go-bluestore/bluestore/bluefs"
	lrdb "github.com/go-bluestore/lib/gorocksdb"
	"syscall"
)

type BlueRocksDirectory struct {
	fs *bluefs.BlueFS
}

func (bd *BlueRocksDirectory) New(fs *bluefs.BlueFS) {
	bd.fs = fs
}

func NewBlueRocksDirectory(fs *bluefs.BlueFS) *BlueRocksDirectory {
	return &BlueRocksDirectory{fs: fs}
}

func NewBlueRocksEnv(fs *bluefs.BlueFS) *BlueRocksEnv {
	return &BlueRocksEnv{
		Wrapper: nil,
		fs:      fs,
	}
}

func NewEnvMirror(a *lrdb.Env, b *lrdb.Env, freeA bool, freeB bool) *BlueRocksEnv {
	//return C.EnvMirror(b, a, freeA, freeB)
	return nil
}

func (br *BlueRocksEnv) CreateDir(dirName string) error {
	return br.fs.MkDir(dirName)
}

func (br *BlueRocksEnv) CreateDirIfMissing(dirName string) error {
	r := br.fs.MkDir(dirName)
	if r != nil && r != syscall.EEXIST {
		return r
	}
	return nil
}

func (br *BlueRocksEnv) DeleteDir(dirName string) error {
	return br.fs.RmDir(dirName)
}

func (br *BlueRocksEnv) NewDirectory(dirName string, result *BlueRocksDirectory) error {
	if !br.fs.DirExists(dirName) {
		return syscall.ENOENT
	}
	result.New(br.fs)
	return nil
}
