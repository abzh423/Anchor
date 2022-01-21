package protocol

import (
	"encoding/binary"
)

func PackInt64Array(arr []int64, bitsPerEntry int) (returnValue []int64) {
	bits := make([]byte, len(arr)*bitsPerEntry)

	for k, value := range arr {
		for b := 0; b < bitsPerEntry; b++ {
			if value&(1<<b) == 1<<b {
				bits[k*bitsPerEntry+b] = 1
			} else {
				bits[k*bitsPerEntry+b] = 0
			}
		}
	}

	for i, l := 0, len(bits); i < l; i += bitsPerEntry {
		data := bits[i : i+bitsPerEntry]

		returnValue = append(returnValue, int64(binary.BigEndian.Uint64(data)))
	}

	return
}
