package types

import "github.com/satori/go.uuid"

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
