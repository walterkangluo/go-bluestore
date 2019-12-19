package types

type BufferList struct {
	data []byte
}

func CreateBufferList() *BufferList {
	return &BufferList{
		data: make([]byte, 0),
	}
}

func (bf *BufferList) Length() uint64 {
	return uint64(len(bf.data))
}

func (bf *BufferList) Encode(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		bf.data = append(bf.data, data[i])
	}
	return bf.data
}
