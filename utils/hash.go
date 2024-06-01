package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashAndTrim(input string) string {
	hash := sha256.Sum256([]byte(input))

	hashStr := hex.EncodeToString(hash[:])

	if len(hashStr) >= 5 {
		return hashStr[len(hashStr)-5:]
	}
	return hashStr
}
