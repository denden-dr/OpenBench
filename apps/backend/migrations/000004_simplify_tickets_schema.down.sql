-- Migration 000004 DOWN: Mark as irreversible
-- This migration cannot be automatically reversed because:
--   1. Satellite tables (ticket_comments, payments, etc.) were archived (renamed), not dropped.
--      Data is still present in their *_archived counterparts.
--   2. Schema changes (ADD COLUMN, ALTER COLUMN) cannot be safely reversed without
--      risking data loss on any rows created after the up migration ran.
--
-- Recovery path:
--   1. Restore from a database backup taken before migration 000004 was applied, OR
--   2. Manually rename archived tables back (e.g. ALTER TABLE ticket_comments_archived RENAME TO ticket_comments)
--      and drop the columns added by the up migration after verifying data integrity.
--
-- Raising an explicit error prevents accidental down-migration attempts.
DO $$
BEGIN
  RAISE EXCEPTION
    'Migration 000004 is intentionally irreversible. '
    'Restore from backup or manually reverse archived table renames. '
    'See migration file header for recovery instructions.';
END;
$$;
