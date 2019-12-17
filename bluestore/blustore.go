package bluestore

import "github.com/go-bluestore/bluestore/common"

type AioContext struct {
}

type TransContext struct {
	AioContext
}

type BlueStore struct {
	common.ObjectStore
	common.BlueFSDeviceExpander
	common.Md_config_obs_t

	ac AioContext
	tc TransContext
}
