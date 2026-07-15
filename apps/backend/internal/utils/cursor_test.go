package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCursorEncodingDecoding(t *testing.T) {
	now := time.Now().UTC()
	id := "test-id-123"

	encoded := EncodeCursor(now, id)
	assert.NotEmpty(t, encoded)

	decodedTime, decodedID, err := DecodeCursor(encoded)
	assert.NoError(t, err)
	assert.Equal(t, id, decodedID)
	assert.True(t, now.Equal(decodedTime), "Times should be equal")
}

func TestCursorInvalidDecoding(t *testing.T) {
	_, _, err := DecodeCursor("invalid-base64-string!!!")
	assert.Error(t, err)

	_, _, err = DecodeCursor("YW55IGNhcm5hbCBwbGVhc3VyZS4=") // decodes to "any carnal pleasure." without "|"
	assert.Error(t, err)
}

func BenchmarkEncodeCursor(b *testing.B) {
	now := time.Now().UTC()
	id := "test-id-123"

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = EncodeCursor(now, id)
	}
}

func BenchmarkDecodeCursor(b *testing.B) {
	now := time.Now().UTC()
	id := "test-id-123"
	encoded := EncodeCursor(now, id)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _, _ = DecodeCursor(encoded)
	}
}
