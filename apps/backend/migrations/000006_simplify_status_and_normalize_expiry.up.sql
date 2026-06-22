-- Update existing tickets: map ready_for_pickup status to completed
UPDATE tickets SET status = 'completed' WHERE status = 'ready_for_pickup';

-- Drop the old constraint
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_status_check;

-- Add the simplified constraint (removing ready_for_pickup)
ALTER TABLE tickets ADD CONSTRAINT tickets_status_check CHECK (status IN ('received', 'in_repair', 'completed', 'cancelled'));

-- Drop the redundant warranty_expiry_date column
ALTER TABLE tickets DROP COLUMN IF EXISTS warranty_expiry_date;
