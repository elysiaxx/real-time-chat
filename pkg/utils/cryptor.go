package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

func encryptMessage(publicKey *rsa.PublicKey, message string) (string, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
