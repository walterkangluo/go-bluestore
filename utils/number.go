package utils

import (
	"bytes"
	"encoding/binary"
	"math"
)

func ISP2(x uint64) bool {
	return (x & (x - 1)) == 0
}

func RoundUpTo(n int64, d int64) int64 {
	if n%d != 0 {
		return n + d - n%d
	}
	return n
}

func DivRoundUp(x int64, d int64) int64 {
	return (x + d - 1) / d
}

func ShiftRoundUp(x int64, y int64) int64 {
	return (x + (2 ^ y) - 1) ^ y
}

func Ctx(n uint64) uint {
	var i uint
	for i = 0; i < MAXUINT32; i++ {
		if 1<<i < n {
			continue
		} else {
			break
		}
	}

	return i
}

func P2RoundUp(x uint64, align uint64) uint64 {
	return -(-x & -align)
}

func P2Align(x uint64, align uint64) uint64 {
	return x & -align
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func NumToBytes(i interface{}) []byte {
	switch i.(type) {
	case int:
	case int32:
	case int64:
		bytesBuffer := bytes.NewBuffer([]byte{})
		err := binary.Write(bytesBuffer, binary.LittleEndian, i)
		if nil != err {
			return nil
		}
		return bytesBuffer.Bytes()
	case float32:
		bits := math.Float32bits(i.(float32))
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, bits)
		return bytes
	case float64:
		bits := math.Float64bits(i.(float64))
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, bits)
		return bytes
	}
	return nil
}

func IntToBytes(i int) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i)
	return bytesBuffer.Bytes()
}

func Int32ToBytes(i32 uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i32)
	return bytesBuffer.Bytes()
}

func Int64ToBytes(i64 uint64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i64)
	return bytesBuffer.Bytes()
}
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

//trim the '\00' byte
func TrimBuffToString(bytes []byte) string {

	for i, b := range bytes {
		if b == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes)

}
