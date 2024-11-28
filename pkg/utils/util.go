package utils

import (
	"crypto/sha256"
	"fmt"
)

func Hash_SHA256(raw string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(raw))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
