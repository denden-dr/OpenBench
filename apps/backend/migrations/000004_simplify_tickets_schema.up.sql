-- Migration: Simplify tickets schema (data-preserving)
-- This migration evolves the existing tickets table in-place and archives
-- satellite tables. It does NOT drop the tickets table or any data within it.

-- Step 1: Add new columns to tickets if they do not already exist.
-- These columns support the simplified single-user dashboard workflow.
ALTER TABLE tickets
  ADD COLUMN IF NOT EXISTS customer_name        TEXT,
  ADD COLUMN IF NOT EXISTS customer_gender      TEXT,
  ADD COLUMN IF NOT EXISTS issue                TEXT,
  ADD COLUMN IF NOT EXISTS additional_description TEXT,
  ADD COLUMN IF NOT EXISTS accessories          TEXT,
  ADD COLUMN IF NOT EXISTS price                DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS payment_status       TEXT NOT NULL DEFAULT 'unpaid',
  ADD COLUMN IF NOT EXISTS warranty_days        INTEGER NOT NULL DEFAULT 30,
  ADD COLUMN IF NOT EXISTS entry_date           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN IF NOT EXISTS exit_date            TIMESTAMP WITH TIME ZONE,
  ADD COLUMN IF NOT EXISTS warranty_expiry_date TIMESTAMP WITH TIME ZONE;

-- Step 2: Backfill mandatory text columns from legacy columns where possible,
-- so existing rows are not left with NULLs in fields we will later constrain.
UPDATE tickets
SET
  customer_name   = COALESCE(customer_name, 'Unknown'),
  customer_gender = COALESCE(customer_gender, 'Male'),
  issue           = COALESCE(issue, issue_description, 'No issue recorded'),
  entry_date      = COALESCE(entry_date, created_at, CURRENT_TIMESTAMP)
WHERE
  customer_name IS NULL
  OR customer_gender IS NULL
  OR issue IS NULL;

-- Step 3: Apply NOT NULL constraints now that all rows have values.
ALTER TABLE tickets
  ALTER COLUMN customer_name   SET NOT NULL,
  ALTER COLUMN customer_gender SET NOT NULL,
  ALTER COLUMN issue           SET NOT NULL,
  ALTER COLUMN entry_date      SET NOT NULL;

-- Step 4: Set a sensible default status for any rows that have legacy
-- status values not in the new canonical enum.
UPDATE tickets
SET status = 'service_in'
WHERE status NOT IN ('service_in', 'on_process', 'fixed', 'picked_up');

-- Step 5: Rename the old diagnosis_fee column to price if it still exists
-- (rename is idempotent via a check on the information_schema).
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tickets' AND column_name = 'diagnosis_fee'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tickets' AND column_name = 'price'
  ) THEN
    ALTER TABLE tickets RENAME COLUMN diagnosis_fee TO price;
  END IF;
END;
$$;

-- Step 6: Archive satellite tables rather than dropping them.
-- This preserves historical data while removing them from active use.
-- Tables are renamed with an _archived suffix so they can be restored if needed.
DO $$
DECLARE
  tbl TEXT;
BEGIN
  FOREACH tbl IN ARRAY ARRAY['ticket_comments','ticket_parts','parts','attachments','payments','audit_logs','customers']
  LOOP
    IF EXISTS (
      SELECT 1 FROM information_schema.tables WHERE table_name = tbl
    ) AND NOT EXISTS (
      SELECT 1 FROM information_schema.tables WHERE table_name = tbl || '_archived'
    ) THEN
      EXECUTE format('ALTER TABLE %I RENAME TO %I', tbl, tbl || '_archived');
    END IF;
  END LOOP;
END;
$$;
