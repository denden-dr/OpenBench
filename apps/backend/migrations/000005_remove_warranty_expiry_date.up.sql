UPDATE tickets
SET exit_date = NULL
WHERE status <> 'picked_up' AND exit_date IS NOT NULL;

ALTER TABLE tickets DROP COLUMN IF EXISTS warranty_expiry_date;