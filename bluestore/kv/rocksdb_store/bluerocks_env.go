package rocksdb_store

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import (
	"github.com/go-bluestore/bluestore/bluefs"
	lrdb "github.com/go-bluestore/lib/rockdb"
)

func NewBlueRocksEnv(fs *bluefs.BlueFS) *BlueRocksEnv {
	return &BlueRocksEnv{
		//Wrapper: nil,
		fs: fs,
	}
}

func NewEnvMirror(a *lrdb.Env, b *lrdb.Env, freeA bool, freeB bool) *BlueRocksEnv {

	//return C.EnvMirror(b, a, freeA, freeB)
	return nil
}

func (br *BlueRocksEnv) CreateDir(dirName string) error {
	// TODO: to be implement
	return br.fs.MkDir(dirName)
}
