-- Migration: Remove warranty_expiry_date column from tickets table
-- Reason: warranty_expiry_date is now computed from exit_date + warranty_days instead of being stored

BEGIN;

CREATE TABLE IF NOT EXISTS _migration_000005_backup AS 
SELECT id, warranty_expiry_date 
FROM tickets 
WHERE warranty_expiry_date IS NOT NULL;

UPDATE tickets
SET exit_date = NULL
WHERE status <> 'picked_up' AND exit_date IS NOT NULL;

ALTER TABLE tickets DROP COLUMN IF EXISTS warranty_expiry_date;

COMMIT;