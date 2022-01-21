package util

import "math/bits"

func CompactInt64Array(array []int64, bitsPerEntry int) (result []int64) {
	result = make([]int64, len(array))

	resultIndex := 0
	bitOffset := 64 - bits.Len64(uint64(array[0]))

	for i := 0; i < len(array); i++ {
		result[resultIndex] = array[i] << int64(bitOffset)

		if i+1 < len(array) {
			nextValueBits := bits.Len64(uint64(array[i+1]))

			if bitOffset+nextValueBits >= 64-1 {
				bitOffset = 64 - nextValueBits
				resultIndex++
			}
		}
	}

	return
}

func CompactInt32Array(array []int32, bitsPerEntry int) (result []int64) {
	result = make([]int64, len(array))

	resultIndex := 0
	bitOffset := 64 - bits.Len32(uint32(array[0]))

	for i := 0; i < len(array); i++ {
		result[resultIndex] = int64(array[i]) << int64(bitOffset)

		if i+1 < len(array) {
			nextValueBits := bits.Len32(uint32(array[i+1]))

			if bitOffset+nextValueBits >= 64-1 {
				bitOffset = 64 - nextValueBits
				resultIndex++
			}
		}
	}

	return
}
