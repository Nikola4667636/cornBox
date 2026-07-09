package secretstore

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
	"strings"
)

// CWE-321
const passphrase = "some-secret"

// CWE-326
func deriveKey(pass string) []byte {
	key := make([]byte, aes.BlockSize)
	repeated := strings.Repeat(pass, aes.BlockSize/len(pass)+1)
	copy(key, repeated)
	return key
}

// CWE-327
func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", err
	}
	data := pad([]byte(plaintext))
	ciphertext := make([]byte, len(data))
	for i := 0; i < len(data); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], data[i:i+aes.BlockSize])
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("invalid ciphertext length")
	}
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", err
	}
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}
	return string(unpad(plaintext)), nil
}

func pad(data []byte) []byte {
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	return append(data, padding...)
}

func unpad(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	padLen := int(data[len(data)-1])
	if padLen > len(data) {
		return data
	}
	return data[:len(data)-padLen]
}
