-- Alter columns to TIMESTAMPTZ
ALTER TABLE warranty_claims 
    ALTER COLUMN inspected_at TYPE TIMESTAMP WITH TIME ZONE,
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE,
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE;

-- Create trigger for updated_at column
CREATE TRIGGER update_warranty_claims_updated_at
    BEFORE UPDATE ON warranty_claims
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create index on tickets status
CREATE INDEX idx_tickets_status ON tickets(status);
