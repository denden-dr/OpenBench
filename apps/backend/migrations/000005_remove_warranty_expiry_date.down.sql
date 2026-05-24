ALTER TABLE tickets ADD COLUMN IF NOT EXISTS warranty_expiry_date TIMESTAMP WITH TIME ZONE;

UPDATE tickets
SET warranty_expiry_date = exit_date + (warranty_days || ' days')::INTERVAL
WHERE status = 'picked_up' AND exit_date IS NOT NULL;
