package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func EncryptAES(data string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Create a new CBC mode cipher
	cipherText := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, []byte(data))

	// Combine IV and ciphertext and base64-encode it
	encryptedData := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func DecryptAES(encryptedData string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract the IV
	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedData := make([]byte, len(cipherText))
	mode.CryptBlocks(decryptedData, cipherText)

	return string(decryptedData), nil
}
