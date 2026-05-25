package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/jmoiron/sqlx"
)

var ErrIdempotencyConflict = errors.New("idempotency key reused with different request body")

type PostgresStorage struct {
	db              *sqlx.DB
	quit            chan struct{}
	closeOnce       sync.Once
	wg              sync.WaitGroup
	cleanupInterval time.Duration
}

type StorageOption func(*PostgresStorage)

func WithCleanupInterval(d time.Duration) StorageOption {
	return func(s *PostgresStorage) {
		s.cleanupInterval = d
	}
}

func NewPostgresStorage(db *sqlx.DB, opts ...StorageOption) *PostgresStorage {
	s := &PostgresStorage{
		db:              db,
		quit:            make(chan struct{}),
		cleanupInterval: 5 * time.Minute,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.cleanupInterval > 0 {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			ticker := time.NewTicker(s.cleanupInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					select {
					case <-s.quit:
						return
					default:
					}
					if err := s.DeleteExpired(); err != nil {
						slog.Error("failed to run background cleanup for idempotency keys", "error", err)
					}
				case <-s.quit:
					return
				}
			}
		}()
	}

	return s
}

func (s *PostgresStorage) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}

	var val []byte
	err := s.db.Get(&val, `
		SELECT value
		FROM idempotency_keys
		WHERE key = $1
		  AND expires_at > CURRENT_TIMESTAMP
		  AND value IS NOT NULL
	`, key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return utils.CopyBytes(val), nil
}

func (s *PostgresStorage) ReserveRequest(key string, requestHash string, exp time.Duration) error {
	if key == "" || requestHash == "" {
		return nil
	}

	expiresAt := time.Now().Add(exp)
	if exp == 0 {
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	var existingHash string
	err := s.db.Get(&existingHash, `
		SELECT request_hash
		FROM idempotency_keys
		WHERE key = $1 AND expires_at > CURRENT_TIMESTAMP
	`, key)
	if err == nil {
		if existingHash != requestHash {
			return ErrIdempotencyConflict
		}
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	result, err := s.db.Exec(`
		INSERT INTO idempotency_keys (key, request_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE
		SET request_hash = EXCLUDED.request_hash,
		    value = NULL,
		    expires_at = EXCLUDED.expires_at
		WHERE idempotency_keys.expires_at <= CURRENT_TIMESTAMP
	`, key, requestHash, expiresAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		err = s.db.Get(&existingHash, `
			SELECT request_hash
			FROM idempotency_keys
			WHERE key = $1 AND expires_at > CURRENT_TIMESTAMP
		`, key)
		if err != nil {
			return err
		}
		if existingHash != requestHash {
			return ErrIdempotencyConflict
		}
	}

	return nil
}

func (s *PostgresStorage) Set(key string, val []byte, exp time.Duration) error {
	if key == "" || len(val) == 0 {
		return nil
	}

	expiresAt := time.Now().Add(exp)
	if exp == 0 {
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	_, err := s.db.Exec(`
		INSERT INTO idempotency_keys (key, request_hash, value, expires_at)
		VALUES ($1, '', $2, $3)
		ON CONFLICT (key) DO UPDATE
		SET value = EXCLUDED.value,
		    expires_at = EXCLUDED.expires_at
	`, key, utils.CopyBytes(val), expiresAt)
	if err != nil {
		slog.Error("failed to write idempotency response cache", "error", err)
	}
	return nil
}

func (s *PostgresStorage) Delete(key string) error {
	if key == "" {
		return nil
	}

	_, err := s.db.Exec("DELETE FROM idempotency_keys WHERE key = $1", key)
	return err
}

func (s *PostgresStorage) Reset() error {
	_, err := s.db.Exec("TRUNCATE TABLE idempotency_keys")
	return err
}

func (s *PostgresStorage) Close() error {
	s.closeOnce.Do(func() {
		close(s.quit)
		s.wg.Wait()
	})
	return nil
}

func (s *PostgresStorage) DeleteExpired() error {
	_, err := s.db.Exec("DELETE FROM idempotency_keys WHERE expires_at <= CURRENT_TIMESTAMP")
	return err
}
