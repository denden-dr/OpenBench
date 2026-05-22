-- Drop unused tables
DROP TABLE IF EXISTS ticket_comments CASCADE;
DROP TABLE IF EXISTS ticket_parts CASCADE;
DROP TABLE IF EXISTS parts CASCADE;
DROP TABLE IF EXISTS attachments CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS customers CASCADE;

-- Recreate tickets table matching the Google Sheet columns
DROP TABLE IF EXISTS tickets CASCADE;

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_name TEXT NOT NULL,
    customer_gender TEXT NOT NULL,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    issue TEXT NOT NULL,
    additional_description TEXT,
    accessories TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    status TEXT NOT NULL DEFAULT 'service_in',
    payment_status TEXT NOT NULL DEFAULT 'unpaid',
    warranty_days INTEGER NOT NULL DEFAULT 30,
    entry_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    exit_date TIMESTAMP WITH TIME ZONE,
    warranty_expiry_date TIMESTAMP WITH TIME ZONE
);
