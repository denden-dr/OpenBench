CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_number VARCHAR(50) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    brand_phone VARCHAR(100) NOT NULL,
    model_phone VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL DEFAULT '',
    damage_description TEXT NOT NULL,
    repair_action TEXT NOT NULL DEFAULT '',
    cost NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'received' CHECK (status IN ('received', 'diagnosing', 'in_repair', 'ready_for_pickup', 'picked_up', 'cancelled')),
    device_position VARCHAR(50) NOT NULL DEFAULT 'warehouse' CHECK (device_position IN ('warehouse', 'picked_up')),
    payment_status VARCHAR(50) NOT NULL DEFAULT 'none' CHECK (payment_status IN ('none', 'requesting', 'paid')),
    payment_method VARCHAR(50) CHECK (payment_method IN ('cash', 'qris')),
    warranty_duration_days INTEGER NOT NULL DEFAULT 30 CHECK (warranty_duration_days >= 0),
    picked_up_at TIMESTAMP WITH TIME ZONE,
    warranty_expiry_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tickets_ticket_number ON tickets(ticket_number);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
