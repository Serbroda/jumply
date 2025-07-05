package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"hash/fnv"
)

func Hash(input string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(input))
	return h.Sum32()
}

func GenerateID(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))[:10]
}
