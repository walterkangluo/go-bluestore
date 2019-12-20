package types

import (
	"testing"
)

func TestCreateUuidD(t *testing.T) {
	u := GenerateRandomUuid()
	t.Log(u)
}
