package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

// EncodeCursor encodes a created_at time and a string ID into an opaque base64 cursor.
func EncodeCursor(t time.Time, id string) string {
	str := fmt.Sprintf("%s|%s", t.Format(time.RFC3339Nano), id)
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// DecodeCursor decodes a base64 cursor back into created_at time and string ID.
func DecodeCursor(encodedCursor string) (time.Time, string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, "", errors.New("invalid cursor format")
	}

	parts := strings.SplitN(string(bytes), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, "", errors.New("invalid cursor values")
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("invalid cursor timestamp: %w", err)
	}

	return t, parts[1], nil
}
