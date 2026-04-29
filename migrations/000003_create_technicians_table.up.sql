-- Up Migration

-- Create the technicians profile table
CREATE TABLE IF NOT EXISTS technicians (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    specialties TEXT,
    rating DECIMAL(3,2) NOT NULL DEFAULT 5.0,
    total_repairs_completed INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Alter bookings to support technician assignments
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS technician_id UUID REFERENCES technicians(user_id) ON DELETE SET NULL;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS estimated_repair_time VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_bookings_technician_id ON bookings(technician_id);
