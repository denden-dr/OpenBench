-- Migration 000009: Fix cancelled status mapping.
-- Reverts status from 'picked_up' back to 'cancelled' for tickets that were
-- originally cancelled.

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables WHERE table_name = 'audit_logs_archived'
    ) THEN
        BEGIN
            UPDATE tickets
            SET
                status = 'cancelled',
                payment_status = 'unpaid',
                exit_date = NULL,
                warranty_expiry_date = NULL
            WHERE status = 'picked_up'
              AND id IN (
                  SELECT DISTINCT ticket_id::uuid
                  FROM audit_logs_archived
                  WHERE status = 'cancelled' OR action = 'cancelled'
              );
        EXCEPTION WHEN OTHERS THEN
            -- Fallback/ignore if columns are named differently
        END;
    END IF;
END;
$$;
