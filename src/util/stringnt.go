package util

func StringNT(value string) []byte {
	bytes := []byte(value)
	bytes = append(bytes, 0x00)

	return bytes
}
