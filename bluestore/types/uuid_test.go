package types

import (
	"testing"
)

func TestCreateUuidD(t *testing.T) {
	u := GenerateRandomUuid()
	t.Log(u)
}

func TestGenerateRandomUuid(t *testing.T) {
}
