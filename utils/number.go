package utils

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
