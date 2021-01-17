package go_fuzz_test

func FuzzTest(data []byte) bool {
	if len(data) == 3 {
		if data[0] == 'F' &&
			data[1] == 'U' &&
			data[2] == 'Z' &&
			data[3] == 'Z' {
			return true
		}
	}
	return false
}

func Fuzz(data []byte) int {
	FuzzTest(data)
	return 0
}
