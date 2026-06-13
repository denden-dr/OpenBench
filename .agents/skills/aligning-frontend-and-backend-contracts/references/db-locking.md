# Concurrency and Row Locking

When executing checks on mutable states (e.g., checking if a token `is_used` before marking it used), a race condition can occur if concurrent requests query the database simultaneously.

### 1. Pessimistic Locking
Use `SELECT ... FOR UPDATE` inside a database transaction to lock the queried row and serialize requests.

```go
// Go Pessimistic Locking Example (RTR / State Check)
tx, err := db.BeginTxx(ctx, nil)
var record TokenRecord
query := `SELECT id, is_used FROM tokens WHERE token_hash = $1 FOR UPDATE`
err = tx.GetContext(ctx, &record, query, tokenHash)
```

### 2. Optimistic Locking
Perform atomic status transitions in the `UPDATE` query directly (e.g., `UPDATE refresh_tokens SET is_used = true WHERE token_hash = $1 AND is_used = false`). Check the affected rows count; if 0, throw a concurrency or reuse conflict.
