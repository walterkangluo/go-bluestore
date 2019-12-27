package bluerocks_env

import "github.com/go-bluestore/bluestore/bluefs"

type BlueRocksEnv struct {
	fs *bluefs.BlueFS
}

func CreateBlueRocksEnv(fs *bluefs.BlueFS) *BlueRocksEnv {
	return &BlueRocksEnv{fs: fs}
}
