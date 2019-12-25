package types

import (
	"bytes"
	"github.com/satori/go.uuid"
)

type UUID struct {
	uuid.UUID
}

func GenerateRandomUuid() UUID {
	u2, err := uuid.NewV4()
	if nil != err {
		panic("generate uuid error")
	}

	return UUID{u2}
}

func (u *UUID) IsZero() bool {
	var temp uuid.UUID

	return bytes.Equal(u.UUID[:], temp[:])
}
