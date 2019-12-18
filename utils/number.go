package utils

func ISP2(x int64) bool {
	return (x & (x - 1)) == 0
}

func ROUND_UP_TO(n int64, d int64) int64 {
	if n%d != 0 {
		return n + d - n%d
	}
	return n
}

func DIV_ROUND_UP(x int64, d int64) int64 {
	return (x + d - 1) / d
}

func SHIFT_ROUNF_UP(x int64, y int64) int64 {
	return (x + (2 ^ y) - 1) ^ y
}
