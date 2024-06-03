package utils

import (
	"crypto/rand"
	"easynote/logs"
	"encoding/base64"
)

func SecureRandString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SimpleRandString(n int) string {
	res, err := SecureRandString(n)
	if err != nil {
		logs.Errorf("[SimpleRandString] SecureRandString err: %v", err)
		return ""
	}
	return res
}
