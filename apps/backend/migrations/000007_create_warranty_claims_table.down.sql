DROP TABLE IF EXISTS warranty_claims;
ALTER TABLE tickets DROP COLUMN IF EXISTS parent_ticket_id;
ALTER TABLE tickets DROP COLUMN IF EXISTS is_warranty;
