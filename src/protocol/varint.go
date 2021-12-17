package protocol

import (
	"bytes"
	"io"
)

func ReadVarInt(r io.Reader) (int32, int, error) {
	var numRead int = 0
	var result int32 = 0

	for {
		data := make([]byte, 1)

		n, err := r.Read(data)

		if err != nil {
			return 0, numRead, err
		}

		if n < 1 {
			return 0, numRead, io.EOF
		}

		value := (data[0] & 0b01111111)
		result |= int32(value) << (7 * numRead)

		numRead++

		if numRead > 5 {
			return 0, numRead, bytes.ErrTooLarge
		}

		if (data[0] & 0b10000000) == 0 {
			break
		}
	}

	return result, numRead, nil
}

func WriteVarInt(value int32, w io.Writer) error {
	var numWritten int = 0

	for {
		if (uint32(value) & 0xFFFFFF80) == 0 {
			_, err := w.Write([]byte{byte(value)})

			numWritten++

			return err
		}

		_, err := w.Write([]byte{byte(value&0x7F | 0x80)})

		if err != nil {
			return err
		}

		value = int32(uint32(value) >> 7)
	}
}
