CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_type TEXT NOT NULL,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    issue_description TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'received',
    diagnosis_fee DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
