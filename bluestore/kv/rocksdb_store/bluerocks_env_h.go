package rocksdb_store

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import (
	"github.com/go-bluestore/bluestore/bluefs"
	lrdb "github.com/go-bluestore/lib/rockdb"
)

type BlueRocksEnv struct {
	Wrapper *lrdb.Env
	fs      *bluefs.BlueFS
}
