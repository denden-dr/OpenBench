ALTER TABLE tickets ADD COLUMN is_warranty BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE tickets ADD COLUMN parent_ticket_id UUID REFERENCES tickets(id) ON DELETE SET NULL;

CREATE TABLE warranty_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    claim_ticket_id UUID REFERENCES tickets(id) ON DELETE SET NULL,
    issue TEXT NOT NULL,
    additional_description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting_inspection',
    void_reason TEXT,
    inspected_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warranty_claims_ticket_id ON warranty_claims(ticket_id);
CREATE INDEX idx_warranty_claims_status ON warranty_claims(status);
