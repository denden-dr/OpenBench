package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubRollbacker struct {
	called bool
	err    error
}

func (s *stubRollbacker) Rollback() error {
	s.called = true
	return s.err
}

func TestRollbackTx_IgnoresRollbackError(t *testing.T) {
	tx := &stubRollbacker{err: errors.New("rollback failed")}

	assert.NotPanics(t, func() {
		rollbackTx(tx)
	})
	assert.True(t, tx.called)
}
