package types

import (
	"bytes"
	"github.com/satori/go.uuid"
)

type UuidD struct {
	uuid.UUID
}

func GenerateRandomUuid() UuidD {
	u2, err := uuid.NewV4()
	if nil != err {
		panic("generate uuid error")
	}

	return UuidD{u2}
}

func (u *UuidD) IsZero() bool {
	var temp uuid.UUID

	return bytes.Equal(u.UUID[:], temp[:])
}