---
name: openbench-inventory-transactions
description: Use when implementing inventory operations, stock deduction, parts logging on tickets, or any database transaction that involves multiple table updates in the OpenBench Go backend
---

# OpenBench Inventory Transaction Patterns

## Overview

Inventory operations (stock deduction, parts logging) **must** happen within a database transaction alongside ticket updates and audit logs. This prevents partial updates that leave data inconsistent.

## The Transaction Pattern

Every inventory operation follows this sequence **atomically**:

```
BEGIN TRANSACTION
  1. SELECT stock_level FROM parts WHERE id = ? FOR UPDATE  (row lock)
  2. Validate: stock_level >= requested quantity
  3. UPDATE parts SET stock_level = stock_level - ?
  4. INSERT INTO ticket_parts (ticket_id, part_id, quantity)
  5. UPDATE tickets SET status = 'repairing' (if applicable)
  6. INSERT INTO audit_logs (action, changes)
COMMIT
```

If **any step fails** → `ROLLBACK` entire transaction.

## Implementation

```go
func (s *ticketService) LogPartsUsed(ctx context.Context, techID, ticketID string, parts []PartUsage) error {
    return s.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        // Verify technician owns this ticket
        ticket, err := s.ticketRepo.FindByID(ctx, tx, ticketID)
        if err != nil { return err }
        if ticket.TechnicianID != techID {
            return ErrForbidden
        }

        for _, p := range parts {
            // 1. Lock row and check stock
            part, err := s.partRepo.FindByIDForUpdate(ctx, tx, p.PartID)
            if err != nil { return err }
            if part.StockLevel < p.Quantity {
                return fmt.Errorf("%w: %s (available: %d, requested: %d)",
                    ErrInsufficientStock, part.Name, part.StockLevel, p.Quantity)
            }

            // 2. Deduct stock
            if err := s.partRepo.DeductStock(ctx, tx, p.PartID, p.Quantity); err != nil {
                return err
            }

            // 3. Link part to ticket
            if err := s.ticketPartRepo.Create(ctx, tx, ticketID, p.PartID, p.Quantity); err != nil {
                return err
            }

            // 4. Audit log
            s.audit.Log(ctx, tx, "deduct_stock", "part", p.PartID, map[string]interface{}{
                "old_stock": part.StockLevel,
                "new_stock": part.StockLevel - p.Quantity,
                "ticket_id": ticketID,
            })
        }

        // 5. Update ticket status
        return s.ticketRepo.UpdateStatus(ctx, tx, ticketID, "repairing")
    })
}
```

## Repository: SELECT FOR UPDATE

```go
func (r *partRepo) FindByIDForUpdate(ctx context.Context, tx *sql.Tx, id string) (*Part, error) {
    var p Part
    err := tx.QueryRowContext(ctx,
        `SELECT id, name, grade, price, stock_level FROM parts WHERE id = $1 FOR UPDATE`,
        id,
    ).Scan(&p.ID, &p.Name, &p.Grade, &p.Price, &p.StockLevel)
    return &p, err
}
```

**`FOR UPDATE`** locks the row until transaction completes, preventing concurrent reads from seeing stale stock.

## Transaction Helper

```go
func (db *DB) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
    }()
    if err := fn(tx); err != nil {
        _ = tx.Rollback()
        return err
    }
    return tx.Commit()
}
```

## Warranty Calculation (on Completion)

```go
func calculateWarrantyExpiry(completedAt time.Time, parts []TicketPart) time.Time {
    minDays := 30 // default to Original (30 days)
    for _, p := range parts {
        if p.Grade == "ODM" {
            minDays = 7 // ODM = 7 days, always wins as minimum
            break
        }
    }
    return completedAt.AddDate(0, 0, minDays)
}
```

## Quick Reference

| Operation | Must Be In Transaction | Needs FOR UPDATE |
|-----------|:---------------------:|:----------------:|
| Deduct stock | ✅ | ✅ |
| Log parts on ticket | ✅ | ❌ |
| Update ticket status | ✅ (with stock ops) | ❌ |
| Audit log insert | ✅ (same tx) | ❌ |
| Payment + ticket completion | ✅ | ❌ |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Checking stock then deducting in separate queries (no tx) | Use `SELECT ... FOR UPDATE` within transaction |
| Forgetting audit log in stock operations | Always insert audit log in same transaction |
| Deducting stock before validation | Validate ALL parts have sufficient stock before any deduction |
| Not using `FOR UPDATE` on parts row | Concurrent requests can oversell without row lock |
| Transaction per part instead of per operation | One transaction for ALL parts in a single request |
| Panicking inside transaction without recover | Use `defer recover` in transaction helper |
