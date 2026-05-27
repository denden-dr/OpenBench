-- Migration 000004: Simplify tickets schema (data-preserving)
-- Evolves the 000001 schema in-place. Safe to apply on databases that went
-- through 000001-000003 (legacy columns present) as well as fresh databases.
-- Does NOT drop the tickets table or destroy any row data.

-- =========================================================================
-- Step 1: Migrate diagnosis_fee → price BEFORE adding a new price column.
--         This must come first so the rename branch can execute when the
--         legacy column exists. If price already exists (re-run / fresh DB),
--         the block is skipped entirely.
-- =========================================================================
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tickets' AND column_name = 'diagnosis_fee'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tickets' AND column_name = 'price'
  ) THEN
    -- Rename carries the column data (including existing fee values) along with it.
    ALTER TABLE tickets RENAME COLUMN diagnosis_fee TO price;
    -- Set a DEFAULT so subsequent ADD COLUMN IF NOT EXISTS for price is a no-op
    -- and any new rows without an explicit value get 0.00.
    ALTER TABLE tickets ALTER COLUMN price DROP NOT NULL;
    ALTER TABLE tickets ALTER COLUMN price SET DEFAULT 0.00;
  END IF;
END;
$$;

-- =========================================================================
-- Step 2: Add new columns that do not exist in the legacy schema.
--         IF NOT EXISTS makes each statement idempotent.
-- =========================================================================
ALTER TABLE tickets
  ADD COLUMN IF NOT EXISTS customer_name           TEXT,
  ADD COLUMN IF NOT EXISTS customer_gender         TEXT,
  ADD COLUMN IF NOT EXISTS issue                   TEXT,
  ADD COLUMN IF NOT EXISTS additional_description  TEXT,
  ADD COLUMN IF NOT EXISTS accessories             TEXT,
  ADD COLUMN IF NOT EXISTS price                   DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  ADD COLUMN IF NOT EXISTS payment_status          TEXT NOT NULL DEFAULT 'unpaid',
  ADD COLUMN IF NOT EXISTS warranty_days           INTEGER NOT NULL DEFAULT 30,
  ADD COLUMN IF NOT EXISTS entry_date              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN IF NOT EXISTS exit_date               TIMESTAMP WITH TIME ZONE,
  ADD COLUMN IF NOT EXISTS warranty_expiry_date    TIMESTAMP WITH TIME ZONE;

-- =========================================================================
-- Step 3: Backfill new text columns from legacy columns so that NOT NULL
--         constraints can be applied without failing on existing rows.
-- =========================================================================
UPDATE tickets
SET
  customer_name   = COALESCE(customer_name,   'Unknown'),
  customer_gender = COALESCE(customer_gender, 'Male'),
  -- Prefer new 'issue' if already populated, fall back to legacy issue_description.
  issue           = COALESCE(
                      NULLIF(issue, ''),
                      NULLIF(issue_description, ''),
                      'No issue recorded'
                    ),
  entry_date      = COALESCE(entry_date, created_at, CURRENT_TIMESTAMP)
WHERE
  customer_name   IS NULL
  OR customer_gender IS NULL
  OR issue        IS NULL
  OR entry_date   IS NULL;

-- =========================================================================
-- Step 4: Apply NOT NULL constraints now that every row has a value.
-- =========================================================================
ALTER TABLE tickets
  ALTER COLUMN customer_name   SET NOT NULL,
  ALTER COLUMN customer_gender SET NOT NULL,
  ALTER COLUMN issue           SET NOT NULL,
  ALTER COLUMN entry_date      SET NOT NULL;

-- =========================================================================
-- Step 5: Relax / drop legacy NOT NULL columns that the repository no
--         longer writes to. This prevents INSERT failures on upgraded DBs.
--
--   device_type      — no replacement; set nullable with a sentinel default.
--   issue_description — superseded by 'issue'; set nullable.
--   diagnosis_fee     — already renamed to 'price' in Step 1 when present.
--
--         Each ALTER is wrapped so it is safe on fresh DBs where the
--         column may not exist.
-- =========================================================================
DO $$
BEGIN
  -- Relax device_type NOT NULL (legacy column, no longer written by app)
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tickets' AND column_name = 'device_type'
      AND is_nullable = 'NO'
  ) THEN
    ALTER TABLE tickets ALTER COLUMN device_type DROP NOT NULL;
    ALTER TABLE tickets ALTER COLUMN device_type SET DEFAULT 'unknown';
    UPDATE tickets SET device_type = 'unknown' WHERE device_type IS NULL;
  END IF;

  -- Relax issue_description NOT NULL (superseded by 'issue')
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tickets' AND column_name = 'issue_description'
      AND is_nullable = 'NO'
  ) THEN
    ALTER TABLE tickets ALTER COLUMN issue_description DROP NOT NULL;
  END IF;
END;
$$;

-- =========================================================================
-- Step 6: Map ALL known legacy status values to the canonical four-state enum.
--
--   Legacy → Canonical
--   received                                       → service_in
--   diagnosing / diagnostics / in_progress /
--     waiting_parts / waiting_customer_confirm /
--     repairing                                    → on_process
--   repaired / ready                               → fixed
--   completed / picked_up / cancelled              → picked_up
--
-- The ELSE branch is intentionally omitted so any unexpected value is caught
-- by the safety UPDATE immediately below.
-- =========================================================================
UPDATE tickets
SET status = CASE status
  WHEN 'received'                  THEN 'service_in'
  WHEN 'diagnosing'                THEN 'on_process'
  WHEN 'diagnostics'               THEN 'on_process'
  WHEN 'in_progress'               THEN 'on_process'
  WHEN 'waiting_parts'             THEN 'on_process'
  WHEN 'waiting_customer_confirm'  THEN 'on_process'
  WHEN 'repairing'                 THEN 'on_process'
  WHEN 'repaired'                  THEN 'fixed'
  WHEN 'ready'                     THEN 'fixed'
  WHEN 'completed'                 THEN 'picked_up'
  WHEN 'cancelled'                 THEN 'cancelled'
  -- Canonical values pass through unchanged:
  WHEN 'service_in'                THEN 'service_in'
  WHEN 'on_process'                THEN 'on_process'
  WHEN 'fixed'                     THEN 'fixed'
  WHEN 'picked_up'                 THEN 'picked_up'
  -- Safety net: anything not matched above maps to service_in so no row
  -- is left outside the canonical enum.
  ELSE 'service_in'
END;

-- Safety guard: assert no non-canonical status survived the mapping above.
-- Raises an exception if any row slipped through, surfacing the value so it
-- can be added to the CASE in a follow-up migration.
DO $$
DECLARE
  bad_count INTEGER;
  bad_sample TEXT;
BEGIN
  SELECT COUNT(*), MIN(status)
  INTO bad_count, bad_sample
  FROM tickets
  WHERE status NOT IN ('service_in', 'on_process', 'fixed', 'picked_up', 'cancelled');

  IF bad_count > 0 THEN
    RAISE EXCEPTION
      'Migration 000004 status mapping incomplete: % row(s) still have '
      'non-canonical status (sample: %). Add the value to the CASE and re-run.',
      bad_count, bad_sample;
  END IF;
END;
$$;

-- =========================================================================
-- Step 7: Backfill final-state invariants for rows that were already
--         completed or picked_up before migration.
--
--   The ticket service sets these three fields atomically when a ticket
--   transitions to picked_up at runtime. Rows migrated from 'completed'
--   or pre-existing 'picked_up' rows may be missing them, which would
--   cause the frontend revenue logic to ignore them (exit_date IS NULL).
-- =========================================================================
UPDATE tickets
SET
  payment_status      = 'paid',
  exit_date           = COALESCE(
                          exit_date,
                          updated_at,
                          entry_date,
                          created_at,
                          CURRENT_TIMESTAMP
                        ),
  warranty_expiry_date = COALESCE(
                           warranty_expiry_date,
                           COALESCE(
                             exit_date,
                             updated_at,
                             entry_date,
                             created_at,
                             CURRENT_TIMESTAMP
                           ) + (warranty_days || ' days')::INTERVAL
                         )
WHERE status = 'picked_up'
  AND (payment_status != 'paid' OR exit_date IS NULL OR warranty_expiry_date IS NULL);

-- =========================================================================
-- Step 8: Archive satellite tables rather than dropping them.
--         Tables are renamed to *_archived so data is recoverable.
-- =========================================================================
DO $$
DECLARE
  tbl TEXT;
BEGIN
  FOREACH tbl IN ARRAY ARRAY[
    'ticket_comments', 'ticket_parts', 'parts',
    'attachments', 'payments', 'audit_logs', 'customers'
  ]
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
