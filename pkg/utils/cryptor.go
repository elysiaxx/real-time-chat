package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

func EncryptMessage(publicKey *rsa.PublicKey, message string) (string, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func DecryptMessage(privateKey *rsa.PrivateKey, encryptedMessage string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(decryptedBytes), nil
}
