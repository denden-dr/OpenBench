-- Update existing status values to match the new options
UPDATE tickets SET status = 'completed' WHERE status = 'picked_up';
UPDATE tickets SET status = 'in_repair' WHERE status = 'diagnosing';

-- Drop the old constraint
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_status_check;

-- Add the new constraint
ALTER TABLE tickets ADD CONSTRAINT tickets_status_check CHECK (status IN ('received', 'in_repair', 'ready_for_pickup', 'completed', 'cancelled'));
