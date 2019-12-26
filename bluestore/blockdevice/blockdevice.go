package blockdevice

import (
	"github.com/go-bluestore/bluestore/types"
	"github.com/go-bluestore/common"
	"github.com/go-bluestore/log"
	"os"
	"strings"
	"unsafe"
)

func CreateBlockDevice(cct *types.CephContext, path string, cb AioCallbackT, cbPriv unsafe.Pointer) *BlockDevice {

	var deviceType = "kernel"
	var r error

	symlink, r := os.Readlink(path)
	if r == nil {
		if strings.HasPrefix(symlink, common.SPDKPrefix) {
			deviceType = "ust-nvme"
		}
	}

	// TODO: Add pmem feature

	log.Debug("path %s and type %s.", path, deviceType)

	if "kernel" == deviceType {
		kernel := CreateKernelDevice(cct, path, cb, cbPriv)
		return kernel.BlockDevice
	}

	// defined HAVE_SPDK
	if "ust-nvme" == deviceType {
		nvme := CreateNVMEDevice(cct, path, cb, cbPriv)
		return nvme.BlockDevice
	}

	log.Error("unsupport type %s.", deviceType)
	// TODO: Add abort process
	return nil
}
