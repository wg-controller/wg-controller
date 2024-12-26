package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func NewWireguardKeyPair() (privateKey, publicKey string, err error) {
	private := make([]byte, 32)
	_, err = rand.Read(private)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes for private key: %w", err)
	}

	// Copy private key to fixed size buffer
	var privKey [32]byte
	copy(privKey[:], private)

	// Generate the public key
	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, &privKey)

	// Convert to base64 encoded strings
	privateKey = base64.StdEncoding.EncodeToString(private)
	publicKey = base64.StdEncoding.EncodeToString(pubKey[:])

	return privateKey, publicKey, nil
}

func NewAESKey() (string, error) {
	return GenerateRandomString(32)
}
