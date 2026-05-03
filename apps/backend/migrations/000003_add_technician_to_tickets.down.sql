-- apps/backend/migrations/000003_add_technician_to_tickets.down.sql
ALTER TABLE tickets DROP COLUMN technician_id;
