DROP INDEX IF EXISTS idx_claims_evaluation_status;

ALTER TABLE claims DROP COLUMN IF EXISTS evaluation_status;
ALTER TABLE claims DROP COLUMN IF EXISTS evaluation_notes;
ALTER TABLE warranties DROP COLUMN IF EXISTS notes;
