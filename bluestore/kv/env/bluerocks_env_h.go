package env

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import (
	"github.com/go-bluestore/bluestore/bluefs"
	lrdb "github.com/go-bluestore/lib/gorocksdb"
)

type BlueRocksEnv struct {
	Wrapper *lrdb.Env
	fs      *bluefs.BlueFS
}
