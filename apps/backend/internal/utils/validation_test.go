package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		s := testStruct{
			Name:  "Denden",
			Email: "denden@example.com",
			Age:   20,
		}
		err := ValidateStruct(s)
		assert.NoError(t, err)
	})

	t.Run("invalid struct - missing fields", func(t *testing.T) {
		s := testStruct{
			Name:  "",
			Email: "invalid-email",
			Age:   15,
		}
		err := ValidateStruct(s)
		assert.Error(t, err)
	})
}

func BenchmarkValidateStruct(b *testing.B) {
	s := testStruct{
		Name:  "Denden",
		Email: "denden@example.com",
		Age:   20,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = ValidateStruct(s)
	}
}
