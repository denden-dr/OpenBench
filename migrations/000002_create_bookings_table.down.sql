-- Down Migration
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP TABLE IF EXISTS bookings;
