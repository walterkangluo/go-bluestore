package utils

func AssertTrue(res bool) {
	if false == res {
		panic("error")
	}
	return
}
