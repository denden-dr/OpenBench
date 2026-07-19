CREATE TABLE IF NOT EXISTS warranties (
    id VARCHAR(100) PRIMARY KEY,
    ticket_id UUID NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS claims (
    id VARCHAR(100) PRIMARY KEY,
    claim_number VARCHAR(100) UNIQUE NOT NULL,
    warranty_id VARCHAR(100) NOT NULL REFERENCES warranties(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'RECEIVED',
    issue_description TEXT NOT NULL,
    repair_action TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_warranties_ticket_id ON warranties(ticket_id);
CREATE INDEX IF NOT EXISTS idx_warranties_status ON warranties(status);
CREATE INDEX IF NOT EXISTS idx_claims_warranty_id ON claims(warranty_id);
CREATE INDEX IF NOT EXISTS idx_claims_status ON claims(status);
CREATE INDEX IF NOT EXISTS idx_claims_number ON claims(claim_number);
