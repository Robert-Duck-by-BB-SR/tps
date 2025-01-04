package hash

func abs(a, b byte) byte {
	if a < b {
		return b - a
	}

	return a - b
}

func Encode(left, right []byte) []byte {
	var result []byte

	if len(right) < len(left) {
		tmp := left
		left = right
		right = tmp
	}

	for i := range left {
		result = append(result, abs(left[i], right[i]))
	}

	return result
}
