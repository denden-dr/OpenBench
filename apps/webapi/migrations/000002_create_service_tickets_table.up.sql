CREATE TABLE IF NOT EXISTS service_tickets (
    id UUID PRIMARY KEY,
    ticket_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'RECEIVED',
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    device_brand VARCHAR(100) NOT NULL,
    device_model VARCHAR(100) NOT NULL,
    device_passcode VARCHAR(100),
    issue_description TEXT NOT NULL,
    repair_action TEXT,
    cost BIGINT NOT NULL DEFAULT 0,
    warranty_days INT NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_service_tickets_number ON service_tickets(ticket_number);
CREATE INDEX IF NOT EXISTS idx_service_tickets_status ON service_tickets(status);
