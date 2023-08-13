package utils

import "math"

func MakeRange(start, end int) []int {
	size := math.Abs(float64(end-start)) + 1
	arr := make([]int, int(size))

	if end > start {
		for i := range arr {
			arr[i] = start + i
		}

		return arr
	}

	for i := range arr {
		arr[i] = start - i
	}

	return arr
}
