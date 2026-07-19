package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// Encrypt encrypts plain text string using AES-GCM with a 32-byte key.
// Returns the hex-encoded string containing the 12-byte nonce prepended to the ciphertext.
func Encrypt(plainText string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("encryption key must be exactly 32 bytes/characters long")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

// Decrypt decrypts a hex-encoded cipher text string using AES-GCM with a 32-byte key.
func Decrypt(cipherTextHex string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("decryption key must be exactly 32 bytes/characters long")
	}

	cipherText, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, actualCipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
