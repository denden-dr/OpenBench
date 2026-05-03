-- apps/backend/migrations/000003_add_technician_to_tickets.up.sql
ALTER TABLE tickets ADD COLUMN technician_id UUID;
