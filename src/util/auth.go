package util

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

// Credit: https://github.com/Tnze/go-mc/blob/c724909cd71137ae3ff490b034a20d7a3c80e694/server/auth/auth.go
func AuthDigest(serverID string, sharedSecret, publicKey []byte) string {
	h := sha1.New()
	h.Write([]byte(serverID))
	h.Write(sharedSecret)
	h.Write(publicKey)
	hash := h.Sum(nil)

	negative := (hash[0] & 0x80) == 0x80

	if negative {
		hash = TwosComplement(hash)
	}

	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")

	if negative {
		res = "-" + res
	}

	return res
}

func TwosComplement(p []byte) []byte {
	carry := true

	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])

		if carry {
			carry = p[i] == 0xff

			p[i]++
		}
	}

	return p
}
