-- Down Migration
DROP INDEX IF EXISTS idx_bookings_technician_id;
ALTER TABLE bookings DROP COLUMN IF EXISTS estimated_repair_time;
ALTER TABLE bookings DROP COLUMN IF EXISTS technician_id;
DROP TABLE IF EXISTS technicians;
