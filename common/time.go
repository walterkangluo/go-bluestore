package common

import (
	"time"
)

var second = uint32(time.Second)

type tsT struct {
	tvSec  uint32
	tvNSec uint32
}

type UTimeT struct {
	tv tsT
}

func (ut *UTimeT) IsZero() bool {
	return ut.tv.tvNSec == 0 && ut.tv.tvSec == 0
}

func (ut *UTimeT) Normalize() {
	if ut.tv.tvNSec > second {
		ut.tv.tvSec += ut.tv.tvNSec / second
		ut.tv.tvNSec %= second
	}
}

func (ut *UTimeT) New() {

}
