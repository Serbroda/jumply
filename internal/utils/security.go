package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateID(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))[:10]
}
