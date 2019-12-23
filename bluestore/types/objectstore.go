package types

type ObjectStore struct {
	Cct *CephContext

	Path string
}

func (obs *ObjectStore)ReadMeta(key string, value *string) int {
	return 0
}
