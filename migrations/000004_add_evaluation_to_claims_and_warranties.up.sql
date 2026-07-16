ALTER TABLE claims ADD COLUMN evaluation_status VARCHAR(50) NOT NULL DEFAULT 'PENDING';
ALTER TABLE claims ADD COLUMN evaluation_notes TEXT;
ALTER TABLE warranties ADD COLUMN notes TEXT;

CREATE INDEX IF NOT EXISTS idx_claims_evaluation_status ON claims(evaluation_status);
