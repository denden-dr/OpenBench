-- Add the warranty_expiry_date column back
ALTER TABLE tickets ADD COLUMN warranty_expiry_date TIMESTAMP WITH TIME ZONE;

-- Restore warranty_expiry_date values from warranties table
UPDATE tickets t 
SET warranty_expiry_date = w.end_date 
FROM warranties w 
WHERE t.id = w.ticket_id;

-- Drop the simplified constraint
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_status_check;

-- Add the original constraint back (including ready_for_pickup)
ALTER TABLE tickets ADD CONSTRAINT tickets_status_check CHECK (status IN ('received', 'in_repair', 'ready_for_pickup', 'completed', 'cancelled'));

-- Restore ready_for_pickup status for tickets that are completed but still in the warehouse
UPDATE tickets SET status = 'ready_for_pickup' WHERE status = 'completed' AND device_position = 'warehouse';
