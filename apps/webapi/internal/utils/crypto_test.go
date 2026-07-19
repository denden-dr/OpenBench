package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "this_is_a_secret_key_32_chars_ok"
	plainText := "my-secret-passcode-12345"

	// Encrypt
	encrypted, err := Encrypt(plainText, key)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plainText, encrypted)

	// Decrypt
	decrypted, err := Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plainText, decrypted)

	// Decrypt with invalid key length
	_, err = Decrypt(encrypted, "shortkey")
	assert.Error(t, err)

	// Encrypt with invalid key length
	_, err = Encrypt(plainText, "shortkey")
	assert.Error(t, err)

	// Decrypt invalid hex string
	_, err = Decrypt("invalid-hex-chars-here!", key)
	assert.Error(t, err)

	// Decrypt ciphertext too short
	_, err = Decrypt("00", key)
	assert.Error(t, err)
}

func BenchmarkEncrypt(b *testing.B) {
	key := "this_is_a_secret_key_32_chars_ok"
	plainText := "my-secret-passcode-12345"

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = Encrypt(plainText, key)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	key := "this_is_a_secret_key_32_chars_ok"
	plainText := "my-secret-passcode-12345"
	encrypted, _ := Encrypt(plainText, key)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = Decrypt(encrypted, key)
	}
}
