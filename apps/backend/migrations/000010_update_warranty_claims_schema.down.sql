DROP INDEX IF EXISTS idx_tickets_status;
DROP TRIGGER IF EXISTS update_warranty_claims_updated_at ON warranty_claims;

ALTER TABLE warranty_claims 
    ALTER COLUMN inspected_at TYPE TIMESTAMP,
    ALTER COLUMN created_at TYPE TIMESTAMP,
    ALTER COLUMN updated_at TYPE TIMESTAMP;
