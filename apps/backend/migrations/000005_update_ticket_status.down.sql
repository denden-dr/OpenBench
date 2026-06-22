-- Update existing status values back
UPDATE tickets SET status = 'picked_up' WHERE status = 'completed';

-- Drop the new constraint
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_status_check;

-- Add the old constraint
ALTER TABLE tickets ADD CONSTRAINT tickets_status_check CHECK (status IN ('received', 'diagnosing', 'in_repair', 'ready_for_pickup', 'picked_up', 'cancelled'));
